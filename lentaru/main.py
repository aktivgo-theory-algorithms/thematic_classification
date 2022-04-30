import pyspark.sql.functions as F
import sparknlp
from pyspark.ml.classification import RandomForestClassifier
from pyspark.ml.feature import HashingTF, IDF, StringIndexer, IndexToString
from sparknlp.annotator import *
from sparknlp.base import *

if __name__ == '__main__':
    spark = sparknlp.start()

    main_df = spark.read.csv(
        '../thematic_classification/lenta-ru-news.csv',
        header=True,
        multiLine=True,
        escape="\"")
    main_df.show()
    print(main_df.count())

    filtered_df = main_df \
        .na.drop(subset=["tags"]) \
        .select(["text", "tags"]) \
        .withColumn("text", F.regexp_replace(F.col("text"), "[\n\r]", " ")) \
        .withColumn("text", F.regexp_replace(F.col("text"), ".Rambler Title ", "")) \
        .withColumn("text", F.trim(F.col("text")))
    filtered_df.show()

    print(filtered_df.count())

    count_df = filtered_df \
        .groupBy("tags") \
        .count() \
        .orderBy(F.col("count").desc())
    count_df.show()

    print(count_df.count())

    selected_rows = count_df \
        .select("tags") \
        .limit(5) \
        .collect()
    selected_tags = [row.tags for row in selected_rows]

    print(selected_tags)

    df = filtered_df \
        .filter(F.col("tags").isin(selected_tags))
    df.show()

    print(df.count())

    train_df, test_df = df.randomSplit([0.8, 0.2], seed=1)

    train_df.groupBy("tags") \
        .count() \
        .orderBy(F.col("count").desc()) \
        .show()

    test_df.groupBy("tags") \
        .count() \
        .orderBy(F.col("count").desc()) \
        .show()

    document_assembler = DocumentAssembler() \
        .setInputCol("text") \
        .setOutputCol("document")

    sentence_detector = SentenceDetector() \
        .setInputCols('document') \
        .setOutputCol('sentence')

    tokenizer = Tokenizer() \
        .setInputCols('sentence') \
        .setOutputCol('token')

    stop_words_cleaner = StopWordsCleaner \
        .pretrained('stopwords_ru', 'ru') \
        .setInputCols("token") \
        .setOutputCol("cleanTokens") \
        .setCaseSensitive(False)

    lemmatizer = LemmatizerModel \
        .pretrained("lemma", "ru") \
        .setInputCols("cleanTokens") \
        .setOutputCol("lemma")

    finisher = Finisher() \
        .setInputCols("lemma") \
        .setOutputCols("token_features") \
        .setOutputAsArray(True) \
        .setCleanAnnotations(False)

    hashing_TF = HashingTF(
        inputCol="token_features",
        outputCol="raw_features")

    idf = IDF(
        inputCol="raw_features",
        outputCol="features",
        minDocFreq=5)

    tags_indexer = StringIndexer(
        inputCol="tags",
        outputCol="label")

    ran_forest = RandomForestClassifier(
        numTrees=10)

    tags_to_string_indexer = IndexToString(
        inputCol="label",
        outputCol="article_class")

    pipeline = Pipeline(
        stages=[
            document_assembler,
            sentence_detector,
            tokenizer,
            stop_words_cleaner,
            lemmatizer,
            finisher,
            hashing_TF,
            idf,
            tags_indexer,
            ran_forest,
            tags_to_string_indexer
        ])

    classification_model = ran_forest.fit(train_df)

    from sklearn.metrics import classification_report, accuracy_score

    df_rf = classification_model \
        .transform(test_df) \
        .select("tags", "label", "prediction", "text")
    df_rf_pandas = df_rf.toPandas()

    labels_df = df_rf \
        .select("label", "tags") \
        .distinct() \
        .orderBy("label")
    labels_df.show()

    labels_raw = labels_df.collect()
    labels = [row.tags for row in labels_raw]

    print(labels)

    print(classification_report(
        df_rf_pandas.label, df_rf_pandas.prediction, target_names=labels))

    print(accuracy_score(
        df_rf_pandas.label, df_rf_pandas.prediction))
