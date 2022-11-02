from src.minute.main import Url, getMinute


def daySpeech(date: str):
    start_recode = "1"
    while (start_recode):
        u = Url(3)
        u.from_(date)
        u.until(date)
        u.recordPacking("json")
        u.startRecord(start_recode)
        print(u.getUrl())
        recode = getMinute(u)
        block = recode.iterator()
        while (block.has_next()):
            spc = block.next().speech()
            sentenses = spc.getSentences()
            for s in sentenses:
                yield s
        if (recode.nextRecordPosition() == None):
            break
        start_recode = str(recode.nextRecordPosition())
        print("nextRecordPosition:" + start_recode)
        print()
