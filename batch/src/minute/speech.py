import re
import neologdn


class Speech:
    """
    発言内容を取得する
    """

    def __init__(self, text) -> None:
        self.t = self.__preprocessing(text)

    def __preprocessing(self, text: str) -> str:
        text = re.sub(r"○.*　", "", text)
        text = re.sub(r'\n', '', text)
        text = re.sub(r'\r', '', text)
        text = re.sub(r'\s', '', text)
        text = neologdn.normalize(text)
        text = text.lower()
        return text

    def text(self):
        return self.t

    def getSentences(self):
        sentences = [s+"。" for s in self.t.split('。') if s]
        return sentences


class Recode:
    """
    一度のFetchで取得するJSON情報
    """

    def __init__(self, response) -> None:
        self.response = response

    def numberOfRecords(self) -> int:
        """
        Returns:
            str: 総結果件数
        """
        return self.response["numberOfRecords"]

    def numberOfReturn(self) -> int:
        """
        Returns:
            str: 返戻件数
        """
        return self.response["numberOfReturn"]

    def nextRecordPosition(self) -> int:
        """
        Returns:
            str: 次開始位置
        """
        return self.response["nextRecordPosition"]

    def speechRecord(self):
        """
        Returns:
            list: 発言情報
        """
        return self.response["speechRecord"]

    def speechRecordAt(self, index: int):
        """
        Returns:
            SpeechRecode: 発言情報
        """
        r = self.response["speechRecord"][index]
        return SpeechRecode(r)

    def iterator(self):
        return RecodeIterator(self)


class RecodeIterator:
    index = 0

    def __init__(self, Recode: Recode) -> None:
        self.Recode = Recode

    def has_next(self):
        if self.index < self.Recode.numberOfReturn():
            return True
        else:
            return False

    def next(self):
        sr = self.Recode.speechRecordAt(self.index)
        self.index += 1
        return sr


class SpeechRecode:
    """
    一つの発言
    """

    def __init__(self, speech_recode) -> None:
        self.speech_recode = speech_recode

    def speechID(self) -> str:
        """
        Returns:
            str: 発言ID
        """
        return self.speech_recode["speechID"]

    def issueID(self) -> str:
        """
        Returns:
            str: 会議録ID
        """
        return self.speech_recode["issueID"]

    def imageKind(self) -> str:
        """
        Returns:
            str: イメージ種別（会議録・目次・索引・附録・追録）
        """
        return self.speech_recode["imageKind"]

    def searchObject(self) -> str:
        """
        Returns:
            str: 検索対象箇所（議事冒頭・本文）
        """
        return self.speech_recode["searchObject"]

    def session(self) -> str:
        """
        Returns:
            str: 国会回次
        """
        return self.speech_recode["session"]

    def nameOfHouse(self) -> str:
        """
        Returns:
            str: 院名
        """
        return self.speech_recode["nameOfHouse"]

    def nameOfMeeting(self) -> str:
        """
        Returns:
            str: 会議名
        """
        return self.speech_recode["nameOfMeeting"]

    def issue(self) -> str:
        """
        Returns:
            str: 号数
        """
        return self.speech_recode["issue"]

    def speechID(self) -> str:
        """
        Returns:
            str: 
        """
        return self.speech_recode["speechID"]

    def date(self) -> str:
        """
        Returns:
            str: 開催日付
        """
        return self.speech_recode["date"]

    def closing(self) -> str:
        """
        Returns:
            str: 閉会中フラグ
        """
        return self.speech_recode["closing"]

    def speechOrder(self) -> str:
        """
        Returns:
            str: 発言番号
        """
        return self.speech_recode["speechOrder"]

    def speaker(self) -> str:
        """
        Returns:
            str: 発言者名
        """
        return self.speech_recode["speaker"]

    def speakerYomi(self) -> str:
        """
        Returns:
            str: 発言者よみ
        """
        return self.speech_recode["speakerYomi"]

    def speakerGroup(self) -> str:
        """
        Returns:
            str: 発言者所属会派
        """
        return self.speech_recode["speakerGroup"]

    def speakerPosition(self) -> str:
        """
        Returns:
            str: 発言者肩書き
        """
        return self.speech_recode["speakerPosition"]

    def speakerRole(self) -> str:
        """
        Returns:
            str: 発言者役割
        """
        return self.speech_recode["speakerRole"]

    def speech(self) -> Speech:
        """
        Returns:
            str: 発言
        """
        s = self.speech_recode["speech"]
        return Speech(s)

    def startPage(self) -> str:
        """
        Returns:
            str: 発言が掲載されている開始ページ
        """
        return self.speech_recode["startPage"]

    def speechURL(self) -> str:
        """
        Returns:
            str: 発言URL
        """
        return self.speech_recode["speechURL"]

    def meetingURL(self) -> str:
        """
        Returns:
            str: 会議録テキスト表示画面のURL
        """
        return self.speech_recode["meetingURL"]

    def pdfURL(self) -> str:
        """
        Returns:
            str: 会議録PDF表示画面のURL（※存在する場合のみ）
        """
        return self.speech_recode["pdfURL"]
