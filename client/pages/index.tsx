import Head from "next/head";
import { env } from "process";
import { useEffect, useState } from "react";

export default function Home() {
    const [isEnv, setEnv] = useState(true);
    useEffect(() => {
        setEnv(env.NODE_ENV !== "production");
    }, []);
    return (
        <div>
            <Head>
                <title>Collokai</title>
                <meta
                    name="description"
                    content="国会議事録APIからコロケーションでその日のトピックをざっくり把握する。"
                />
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <div>
                <h1>Collokai</h1>
                <div>
                    国会議事録APIからコロケーションでその日のトピックをざっくり把握する
                </div>
                <div>{isEnv ? "開発" : "本番"}</div>
            </div>
        </div>
    );
}
