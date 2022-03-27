from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.ensemble import RandomForestClassifier
from sklearn.pipeline import Pipeline


def get_thematic(text):
    texts = ['Текст номер один', 'Текст номер два', 'Компьютеры в лингвистике', 'Компьютеры и обработка текстов']
    texts_labels = ['Тема про текст', 'Тема текст', 'Тема про компьютер', 'Тема про компьютер']

    text_clf = Pipeline([
        ('tfidf', TfidfVectorizer()),
        ('clf', RandomForestClassifier())
    ])

    text_clf.fit(texts, texts_labels)

    res = text_clf.predict([text])
    return res


if __name__ == '__main__':
    input_text = input("Введите текст: ")
    result_text = get_thematic(input_text)
    print(result_text)
