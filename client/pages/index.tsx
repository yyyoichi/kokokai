import Head from "next/head";

export default function Home() {
    return (
        <>
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
            </div>
        </>
    );
}
