import re

import nltk
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import confusion_matrix, classification_report, accuracy_score

nltk.download('stopwords')
nltk.download('wordnet')
nltk.download('omw-1.4')

import nltk
import pandas as pd
from matplotlib import pyplot as plt
from sklearn.discriminant_analysis import LinearDiscriminantAnalysis
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import train_test_split, StratifiedKFold, cross_val_score
from sklearn.naive_bayes import GaussianNB
from sklearn.neighbors import KNeighborsClassifier
from sklearn.svm import SVC
from sklearn.tree import DecisionTreeClassifier
from wordcloud import WordCloud, STOPWORDS

nltk.download("stopwords")
from nltk.corpus import stopwords

russian_stopwords = stopwords.words("russian")


# Получение текстовой строки из списка слов
def str_corpus(corpus):
    str_corpus = ''
    for i in corpus:
        str_corpus += ' ' + i
    str_corpus = str_corpus.strip()
    return str_corpus


# Получение списка всех слов в корпусе
def get_corpus(data):
    corpus = []
    for phrase in data:
        for word in phrase.split():
            corpus.append(word)
    return corpus


# Получение облака слов
def get_wordCloud(corpus):
    wordCloud = WordCloud(background_color='white',
                          stopwords=STOPWORDS,
                          width=3000,
                          height=2500,
                          max_words=200,
                          random_state=42
                          ).generate(str_corpus(corpus))
    return wordCloud


# Удаление знаков пунктуации из текста
def remove_punct(text):
    table = {33: ' ', 34: ' ', 35: ' ', 36: ' ', 37: ' ', 38: ' ', 39: ' ', 40: ' ', 41: ' ', 42: ' ', 43: ' ', 44: ' ',
             45: ' ', 46: ' ', 47: ' ', 58: ' ', 59: ' ', 60: ' ', 61: ' ', 62: ' ', 63: ' ', 64: ' ', 91: ' ', 92: ' ',
             93: ' ', 94: ' ', 95: ' ', 96: ' ', 123: ' ', 124: ' ', 125: ' ', 126: ' '}
    return text.translate(table)


if __name__ == '__main__':
    import sys
    import csv

    csv.field_size_limit(sys.maxsize)

    # url = "../../datasetcreator/habrposts.csv"
    # url = "../../habr_posts.csv"
    url = "../../convertcsv.csv"
    names = ['text', 'class']
    dataset = pd.read_csv(url, names=names, engine='python', encoding='utf-8', error_bad_lines=False)
    # print(dataset.head())
    # print(dataset.shape)
    # print(dataset.describe())
    # print(dataset.groupby('tag', sort=False).size())

    dataset['text'] = dataset['text'].map(lambda x: x.replace('\n', ''))

    # Разделение датасета на обучающую и контрольную выборки
    array = dataset.values

    # Выбор первых 4-х столбцов
    X = array[:, 0]

    # Выбор 5-го столбца
    y = array[:, 1]

    documents = []

    from nltk.stem import WordNetLemmatizer

    stemmer = WordNetLemmatizer()

    for sen in range(0, len(X)):
        # Remove all the special characters
        document = re.sub(r'\W', ' ', str(X[sen]))

        # remove all single characters
        document = re.sub(r'\s+[a-zA-Z]\s+', ' ', document)

        # Remove single characters from the start
        document = re.sub(r'\^[a-zA-Z]\s+', ' ', document)

        # Substituting multiple spaces with single space
        document = re.sub(r'\s+', ' ', document, flags=re.I)

        # Removing prefixed 'b'
        document = re.sub(r'^b\s+', '', document)

        # Converting to Lowercase
        document = document.lower()

        # Lemmatization
        document = document.split()

        document = [stemmer.lemmatize(word) for word in document]
        document = ' '.join(document)

        documents.append(document)

    from sklearn.feature_extraction.text import CountVectorizer, TfidfVectorizer

    vectorizer = CountVectorizer(max_features=1500, min_df=2, max_df=0.95, stop_words=russian_stopwords)
    X = vectorizer.fit_transform(documents).toarray()

    # tfidfconverter = TfidfVectorizer(max_features=1500, min_df=5, max_df=0.7, stop_words=russian_stopwords)
    # X = tfidfconverter.fit_transform(documents).toarray()

    # Разделение X и y на обучающую и контрольную выборки
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)

    classifier = RandomForestClassifier()
    classifier.fit(X_train, y_train)

    y_pred = classifier.predict(X_test)

    print(confusion_matrix(y_test, y_pred))
    print(classification_report(y_test, y_pred))
    print(accuracy_score(y_test, y_pred))

    dataset['title'] = dataset['title'].map(lambda x: x.lower())
    dataset['title'] = dataset['title'].map(lambda x: x.replace('\n', ''))
    dataset['title'] = dataset['title'].map(lambda x: remove_punct(x))
    dataset['title'] = dataset['title'].map(lambda x: x.split(' '))
    dataset['title'] = dataset['title'].map(
        lambda x: [token for token in x if token not in russian_stopwords
                   and token != " "
                   and token.strip() not in punctuation])
    dataset['title'] = dataset['title'].map(lambda x: ' '.join(x))

    # Загружаем алгоритмы модели
    # models = []
    # models.append(('LR', LogisticRegression(solver='liblinear', multi_class='ovr')))
    # models.append(('LDA', LinearDiscriminantAnalysis()))
    # models.append(('KNN', KNeighborsClassifier()))
    # models.append(('CART', DecisionTreeClassifier()))
    # models.append(('NB', GaussianNB()))
    # models.append(('SVM', SVC(gamma='auto')))
    #
    # # оцениваем модель на каждой итерации
    # results = []
    # names = []
    #
    # for name, model in models:
    #     kfold = StratifiedKFold(n_splits=5, random_state=1, shuffle=True)
    #     cv_results = cross_val_score(model, X_train, Y_train, cv=kfold, scoring='accuracy')
    #     results.append(cv_results)
    #     names.append(name)
    #     print('%s: %f (%f)' % (name, cv_results.mean(), cv_results.std()))
    #
    # plt.boxplot(results, labels=names)
    # plt.title('Algorithm Comparison')
    # plt.show()

    # sgd_ppl_clf = Pipeline([
    #     ('tfidf', TfidfVectorizer()),
    #     ('sgd_clf', SGDClassifier(random_state=42))
    # ])
    #
    # knb_ppl_clf = Pipeline([
    #     ('tfidf', TfidfVectorizer()),
    #     ('knb_clf', KNeighborsClassifier(n_neighbors=4))
    # ])
    #
    # sgd_ppl_clf.fit(X_train, Y_train)
    #
    # predicted_sgd = sgd_ppl_clf.predict(X_validation)
    # print(metrics.classification_report(predicted_sgd, Y_validation))

    # # Загружаем алгоритмы модели
    # models = []
    # models.append(('LR', LogisticRegression(solver='liblinear', multi_class='ovr')))
    # models.append(('LDA', LinearDiscriminantAnalysis()))
    # models.append(('KNN', KNeighborsClassifier()))
    # models.append(('CART', DecisionTreeClassifier()))
    # models.append(('NB', GaussianNB()))
    # models.append(('SVM', SVC(gamma='auto')))
    #
    # # оцениваем модель на каждой итерации
    # results = []
    # names = []
    #
    # for name, model in models:
    #     kfold = StratifiedKFold(n_splits=5, random_state=1, shuffle=True)
    #     cv_results = cross_val_score(model, X_train, Y_train, cv=kfold, scoring='accuracy')
    #     results.append(cv_results)
    #     names.append(name)
    #     print('%s: %f (%f)' % (name, cv_results.mean(), cv_results.std()))

    # corpus = get_corpus(dataset['text'])
    # procWordCloud = get_wordCloud(corpus)
    # fig = plt.figure(figsize=(20, 8))
    # plt.subplot(1, 2, 1)
    # plt.imshow(procWordCloud)
    # plt.axis('off')
    # plt.subplot(1, 2, 1)
    # plt.show()
