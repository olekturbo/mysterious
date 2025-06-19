import os
import pickle
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.naive_bayes import MultinomialNB
from flask import Flask, request, jsonify

MODEL_PATH = 'spam_model.pkl'
VECTORIZER_PATH = 'vectorizer.pkl'

def train_and_save():
    df = pd.read_csv('spam.csv', encoding='latin-1')[['v1', 'v2']]
    df.columns = ['label', 'text']
    X = df['text']
    y = df['label'].map({'ham': 0, 'spam': 1})
    X_train, _, y_train, _ = train_test_split(X, y, test_size=0.2, random_state=42)
    vectorizer = TfidfVectorizer()
    X_train_vec = vectorizer.fit_transform(X_train)
    model = MultinomialNB()
    model.fit(X_train_vec, y_train)
    pickle.dump(model, open(MODEL_PATH, 'wb'))
    pickle.dump(vectorizer, open(VECTORIZER_PATH, 'wb'))
    return model, vectorizer

if os.path.exists(MODEL_PATH) and os.path.exists(VECTORIZER_PATH):
    model = pickle.load(open(MODEL_PATH, 'rb'))
    vectorizer = pickle.load(open(VECTORIZER_PATH, 'rb'))
else:
    model, vectorizer = train_and_save()

app = Flask(__name__)

@app.route('/predict', methods=['POST'])
def predict():
    data = request.get_json()
    message = data.get('message', '')
    vect = vectorizer.transform([message])
    prediction = model.predict(vect)[0]
    label = 'spam' if prediction == 1 else 'ham'
    return jsonify({'result': label})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)

