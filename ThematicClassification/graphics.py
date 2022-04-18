# Загрузка библиотек

from pandas import read_csv
from pandas.plotting import scatter_matrix
from matplotlib import pyplot

# Загрузка датасета
url = "../datasetcreator/habrposts.csv"
names = ['title', 'text', 'class']
dataset = read_csv(url, names=names, engine='python', encoding='utf-8', error_bad_lines=False)

# Диаграмма размаха
dataset.plot(kind='box', subplots=True, layout=(2, 2), sharex=False, sharey=False)
pyplot.show()

# Гистограмма распределения атрибутов датасета
dataset.hist()
pyplot.show()

# Матрица диаграмм рассеяния
scatter_matrix(dataset)
pyplot.show()
