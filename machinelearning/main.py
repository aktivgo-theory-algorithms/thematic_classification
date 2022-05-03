import pandas
import re
import numpy
from sklearn.ensemble import RandomForestClassifier, AdaBoostClassifier
from sklearn.neighbors import KNeighborsClassifier
from sklearn.naive_bayes import GaussianNB
from sklearn.tree import DecisionTreeClassifier
from sklearn.neural_network import MLPClassifier
from sklearn.discriminant_analysis import QuadraticDiscriminantAnalysis

def read_dataset():
    url = "../lenta-ru-news.csv"
    return pandas.read_csv(url)


if __name__ == '__main__':
    # read dataset
    main_df = read_dataset()

    # drop unnecessary columns and rows
    main_df = main_df.drop(columns=['url', 'title', 'topic', 'date'])
    main_df = main_df[main_df['tags'].notna()]
    main_df = main_df.drop(main_df[main_df.tags == 'Все'].index)
    main_df = main_df.drop(main_df[main_df.tags == '69-я параллель'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Аналитика рынка'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Экология'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Финансы компаний'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Туризм'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Страноведение'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Производители'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Первая мировая'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Нацпроекты'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Наследие'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Выборы'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Вещи'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Вооружение'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Госрегулирование'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Деньги'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Достижения'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Мемы'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Социальная сфера'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Авто'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Автобизнес'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Английский футбол'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Вирусные ролики'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Вкусы'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Внешний вид'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Жизнь'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Казахстан'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Киберпреступность'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Мировой бизнес'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Мировой опыт'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Мнения'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Молдавия'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Офис'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Преступная Россия'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Прибалтика'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Реклама'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Россия'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Софт'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Белоруссия'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Движение'].index)
    main_df = main_df.drop(main_df[main_df.tags == 'Закавказье'].index)

    print(main_df)

    documents = []

    from nltk.stem import WordNetLemmatizer

    stemmer = WordNetLemmatizer()

    for sen in main_df.text:
        # Remove all the special characters
        document = re.sub(r'\W', ' ', str(sen))

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

    y = main_df.tags

    from nltk.corpus import stopwords
    from sklearn.feature_extraction.text import TfidfVectorizer
    tfidfconverter = TfidfVectorizer(max_features=1500, min_df=5, max_df=0.7, stop_words=stopwords.words('russian'))
    X = tfidfconverter.fit_transform(documents).toarray()

    from sklearn.model_selection import train_test_split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=0)

    classificator = MLPClassifier(alpha=0.05, max_iter=1000)
    classificator.fit(X_train, y_train)
    y_pred = classificator.predict(X_test)

    from sklearn.metrics import classification_report, confusion_matrix, accuracy_score

    print(confusion_matrix(y_test,y_pred))
    print(classification_report(y_test,y_pred))
    print(accuracy_score(y_test, y_pred))