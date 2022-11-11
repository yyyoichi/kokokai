import collections
import datetime
import itertools
import sys
from src.sentence import daySpeech
from src.mcb import getMecab, Morpheme
from src.db import get_connection


class Validation:
    stop_words = [""]

    def __init__(self, morpheme: Morpheme) -> None:
        self.m = morpheme

    def is_stop_word(self):
        return self.m.surface() in self.stop_words

    def is_noun(self):
        return self.m.part_of_speech() == "名詞"

    def is_int(self):
        return self.m.part_of_speech_details_1() == "数"

    def is_asterisk(self):
        return self.m.prototype() == "*"


with open('stop_words.txt', 'r', encoding="utf-8") as f:
    Validation.stop_words = [w.strip() for w in f.readlines()]


def get_nouns(date: str):
    """
    ago 日前の議事録から文章ごとの名詞リストを取得する
    """
    m = getMecab()
    speech = daySpeech(date)
    sentence = next(speech, None)
    # 文章ごとの名詞リスト
    noun_list = []

    def is_target(mph: Morpheme):
        v = Validation(mph)
        if v.is_noun() and not v.is_asterisk() and not v.is_int() and not v.is_stop_word():
            return True
        else:
            return False

    while (sentence):
        p = m.parse(sentence)
        morpheme = next(p, None)
        sentence_noun_list = []
        while (morpheme):
            if is_target(morpheme):
                sentence_noun_list.append(morpheme.prototype())
            morpheme = next(p, None)
        noun_list.append(sentence_noun_list)
        sentence = next(speech, None)

    return noun_list


def get_upper(pl: list, minfreq: int, max: int) -> list[tuple[tuple[str], int]]:
    """
    minfreq以上のペアリストを最大max件取得する
    Args
        pl フラットなペアリスト
        minfreq 最低頻度
        max 取得上限
    """
    dcnt = collections.Counter(pl)
    pl = [(k, dcnt[k]) for k in dcnt.keys() if dcnt[k] >= minfreq]
    return sorted(dcnt.items(), key=lambda x: x[1], reverse=True)[:max]


def main(days: int):
    # 10日前よりも最近のデータは取らないぞ
    if days < 10:
        return

    print(datetime.datetime.now(datetime.timezone(
        datetime.timedelta(hours=9))))
    """
    その日の共起リストをDBに格納する
    """
    date = (datetime.datetime.now(datetime.timezone(
        datetime.timedelta(hours=9))) - datetime.timedelta(days=days)).strftime("%Y-%m-%d")

    # 文章ごとの名詞リスト
    noun_list = get_nouns(date)
    # 文章ごとのペアリスト
    double_pair_list = [
        list(itertools.combinations(nl, 2))
        for nl in noun_list if len(nl) >= 2
    ]
    # フラットなペアリスト
    pair_list = []
    for p in double_pair_list:
        pair_list.extend(p)

    upper_list = get_upper(pair_list, 5, 100)

    word_pk_dict = {}

    conn = get_connection()
    cursor = conn.cursor()

    # kyokiday を新規
    cursor.execute("select nextval('kyokiday_pk_seq')")
    kyoki_day_pk = cursor.fetchone()
    if kyoki_day_pk:
        kyoki_day_pk = kyoki_day_pk[0]

    # kyokiday新規挿入
    cursor.execute(
        "insert into kyokiday(pk, date) values(%s, %s)",
        (kyoki_day_pk, date,)
    )
    print("insert kyokiday! pk:", kyoki_day_pk, "date:", date)

    # 共起リストを1つずつ格納
    for pair in upper_list:
        cursor.execute("select nextval('kyoki_pk_seq')")
        kyoki_pk = cursor.fetchone()
        if kyoki_pk:
            kyoki_pk = kyoki_pk[0]
        # kyoki 新規挿入
        cursor.execute(
            "insert into kyoki(pk, kyokiday, freq) values(%s, %s, %s)",
            (kyoki_pk, kyoki_day_pk, pair[1])
        )
        print("\tinsert kyoki! pk:", kyoki_pk, "freq:", pair[1])
        for word in pair[0]:
            word_pk = word_pk_dict.get(word, None)
            if not word_pk:
                # DBから取得するか、新しくワードを挿入して主キーを取得する
                # DBから主キーを取得
                cursor.execute("select code from word where word=%s", (word,))
                db_word_key = cursor.fetchone()
                if db_word_key:
                    word_pk = db_word_key[0]

                # DBにワードがない
                else:
                    # 新規ワードを作成し、主キーを保持する
                    cursor.execute("select nextval('word_code_seq')")
                    word_pk = cursor.fetchone()
                    if word_pk:
                        word_pk = word_pk[0]
                        word_pk_dict[word] = word_pk
                    cursor.execute(
                        "insert into word(code ,word) values(%s, %s)",
                        (word_pk, word,)
                    )
                    print("\t\t\tinsert word! pk:", word_pk)

            # ペア内容（kyokiitem）を挿入
            cursor.execute(
                "insert into kyokiitem(kyokiday, kyoki, word) values(%s, %s, %s)",
                (kyoki_day_pk, kyoki_pk, word_pk)
            )
            print("\t\tinsert kyokiitem! word:", word)
    cursor.close()
    conn.commit()
    conn.close()
    print("end\n\n")


    # conn.commit()
if __name__ == "__main__":
    arg = sys.argv[1]
    if arg is None:
        arg = 10
    main(arg)
