import collections
import datetime
import itertools
from src.sentence import daySpeech
from src.mcb import getMecab, Morpheme


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


def get_nouns_days_ago(ago: int):
    """
    ago 日前の議事録から文章ごとの名詞リストを取得する
    """
    m = getMecab()
    date = datetime.datetime.now(datetime.timezone(
        datetime.timedelta(hours=9))) - datetime.timedelta(days=ago)
    speech = daySpeech(date.strftime("%Y-%m-%d"))
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
    print(len(pl))
    return sorted(dcnt.items(), key=lambda x: x[1], reverse=True)[:max]


def main():
    """
    その日の共起リストをDBに格納する
    """
    # 文章ごとの名詞リスト
    noun_list = get_nouns_days_ago(10)
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
    for ul in upper_list:
        print(ul)


if __name__ == "__main__":
    main()
