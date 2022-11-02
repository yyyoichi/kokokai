from src.main import daySpeech
from src.word import getMecab, Word


class Validation:
    stop_words = [""]

    def __init__(self, word: Word) -> None:
        self.w = word

    def is_stop_word(self):
        return self.w.surface() in self.stop_words

    def is_noun(self):
        return self.w.part_of_speech() == "名詞"

    def is_int(self):
        return self.w.part_of_speech_details_1() == "数"

    def is_asterisk(self):
        return self.w.prototype() == "*"


with open('stop_words.txt', 'r', encoding="utf-8") as f:
    Validation.stop_words = [w.strip() for w in f.readlines()]


def main():
    m = getMecab()
    date = "2022-10-27"
    it = daySpeech(date)
    s = next(it, None)
    while (s):
        p = m.parse(s)
        w = next(p, None)
        while (w):
            w.surface()
            # print(w.surface() + " - " + w.part_of_speech())
            w = next(p, None)
        s = next(it, None)


if __name__ == "__main__":
    main()
