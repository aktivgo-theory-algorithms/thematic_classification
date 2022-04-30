import pandas
from matplotlib import pyplot
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.neighbors import KNeighborsClassifier
from sklearn.model_selection import cross_val_score
from sklearn.model_selection import StratifiedKFold

pandas.set_option('display.max_rows', None)


def read_test_dataset():
    from sklearn.datasets import load_digits
    return load_digits()


def read_dataset():
    url = "../datasetcreator/dataset.csv"
    names = ['text', 'tag']
    return pandas.read_csv(url, names=names, engine='python', encoding='utf-8', error_bad_lines=False)


def get_info(d):
    print(d.shape)
    print(d.head(20))
    print(d.describe())
    print(d.groupby('tag').size().sort_values(ascending=False))


def get_most_popular_tags(d):
    return d.groupby('tag').size().sort_values(ascending=False)[7:12]


def get_dataset_with_popular_tags(d, pt):
    return d[d['tag'] in pt]


if __name__ == '__main__':
    dataset = read_test_dataset()

    from sklearn.model_selection import train_test_split

    X_train, X_test, y_train, y_test = train_test_split(dataset.data, dataset.target, random_state=11)

    from sklearn.neighbors import KNeighborsClassifier

    knn = KNeighborsClassifier()

    knn.fit(X=X_train, y=y_train)

    predicted = knn.predict(X=X_test)
    expected = y_test

    print(predicted[:20])
    print(expected[:20])

    wrong = [(p, e) for (p, e) in zip(predicted, expected) if p != e]
    print(wrong)

    print(f'{knn.score(X_test, y_test):.2%}')

    from sklearn.metrics import confusion_matrix

    confusion = confusion_matrix(y_true=expected, y_pred=predicted)
    print(confusion)

    from sklearn.metrics import classification_report

    names = [str(digit) for digit in dataset.target_names]
    print(classification_report(expected, predicted, target_names=names))

    # get_info(dataset)
    # mostPopularTags = get_most_popular_tags(dataset)
    # newDataset = get_dataset_with_popular_tags(dataset, mostPopularTags[:, 0])
    # print(newDataset)

# array = dataset.values
# Xdirty = array[:, 0]
# Ydirty = array[:, 1]
#
# print(Xdirty.sort_values())
#
# Xclear = []
# Yclear = []
# i = 0
# for row in Xdirty:
#     if isinstance(row, str):
#         Xclear.append(row.replace(',', ''))
#         Yclear.append(Ydirty[i])
#     i += 1
#
# print(Yclear)

# dataset.text = (dataset.text.str.split()).apply(lambda x: float(x[0].replace(',', '')))

# myDataset = dataset[:1000]
#
# print(myDataset.shape())

# X_train, X_test, y_train, y_test = train_test_split(myDataset.text, myDataset.tag, random_state=11, test_size=0.20)
#
# knn = KNeighborsClassifier()
# knn.fit(X=X_train, y=y_train)
#
# predicted = knn.predict(X=X_test)
# expected = y_test
#
# wrong = [(p, e) for (p, e) in zip(predicted, expected) if p != e]
#
# print(wrong)
#
# print(f'{knn.score(X_test, y_test):.2%}')
#
# from sklearn.metrics import confusion_matrix
# confusion = confusion_matrix(y_true=expected, y_pred=predicted)
# print(confusion)
#
# from sklearn.metrics import classification_report
# names = [str(tag) for tag in myDataset.tag]
# print(classification_report(expected, predicted, target_names=names))
