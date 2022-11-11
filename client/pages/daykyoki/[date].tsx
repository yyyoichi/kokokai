import { GetStaticProps } from "next";

type DayKyoki = {
    date: string;
    kyoki: {
        pk: number;
        freq: number;
        words: string[];
    }[];
};

export default function DayKyoki({ date, kyoki }: DayKyoki) {
    console.log(date);
    console.log(kyoki);
    return <div>hello</div>;
}

export async function getStaticPaths() {
    return {
        paths: [],
        fallback: "blocking",
    };
}
export const getStaticProps: GetStaticProps = async (context) => {
    let date = context.params?.date;
    if (!date) {
        return { notFound: true };
    }
    if (typeof date !== "string") {
        [date] = date;
    }
    const isDateFormat =
        /[0-9]{4}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])/g.test(date);
    if (!isDateFormat) {
        return { notFound: true };
    }
    date = date.replace(
        /([0-9]{4})(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])/,
        "$1-$2-$3"
    );
    let json: DayKyoki;
    const host =
        process.env.VERCEL_ENV === "production"
            ? "https://kokokai.uw.r.appspot.com/"
            : process.env.VERCEL_ENV === "dev"
            ? "http://localhost:8080/"
            : "https://kokokaidev-dot-kokokai.uw.r.appspot.com/";
    console.log(process.env.NEXT_PUBLIC_NODE_ENV);
    const url = `${host}daykyoki?d=${date}`;
    const res = await fetch(url).catch((e) => {
        console.log(e);
        throw Error("not found: " + date);
    });
    if (res.status !== 200) {
        throw Error("not found: " + date);
    }
    json = await res.json();
    console.log(process.env.NODE_ENV, url, json.kyoki ? json.kyoki[0] : "null");
    return {
        // Passed to the page component as props
        props: { date: json.date, kyoki: json.kyoki },
    };
};
