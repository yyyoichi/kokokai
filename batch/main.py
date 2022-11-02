from src.main import daySpeech
from src.word import getMecab


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
            print(w.surface() + " - " + w.part_of_speech())
            w = next(p, None)
        s = next(it, None)


if __name__ == "__main__":
    main()
