import { GetStaticProps } from "next";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Alert from "@mui/material/Alert";
import Container from "@mui/material/Container";
import { Network } from "vis-network/standalone/umd/vis-network.min.js";
import { useEffect, useRef } from "react";
import NetworkNode from "../../src/NetworkNode";
import { Box } from "@mui/material";
import Dateupdown from "../../util/Dateupdown";

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
    const graphRef = useRef<HTMLDivElement>(null!);
    const nn = new NetworkNode();
    kyoki.forEach(({ words }) => {
        nn.addNode(words[0]);
        nn.addNode(words[1]);
        nn.addEdge(words[0], words[1]);
    });
    const nodeData = nn.getNodeData();
    useEffect(() => {
        new Network(graphRef.current, nodeData, { height: "450px" });
    }, [nodeData]);
    return (
        <>
            <Box pt={2}>
                <Container maxWidth="sm">
                    <Box border={"solid 1px gray"} height={"450px"}>
                        <div ref={graphRef}></div>
                    </Box>
                    <Dateupdown dateString={date} />
                    {kyoki.length > 0 ? (
                        <Alert severity="info">出現回数: 共起ワード</Alert>
                    ) : (
                        <Alert severity="info">データがありません</Alert>
                    )}
                    {kyoki.length > 0 ? (
                        <>
                            <List>
                                {kyoki.map((x, i) => {
                                    return (
                                        <ListItem key={i} disablePadding>
                                            <ListItemButton>
                                                <ListItemText
                                                    primary={`${x.freq}: ${x.words[0]}, ${x.words[1]}`}
                                                />
                                            </ListItemButton>
                                        </ListItem>
                                    );
                                })}
                            </List>
                        </>
                    ) : (
                        <></>
                    )}
                </Container>
            </Box>
        </>
    );
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
            ? "https://molten-mariner-368507.uw.r.appspot.com/"
            : process.env.VERCEL_ENV === "dev"
            ? "http://localhost:8080/"
            : "https://collokaidev-dot-molten-mariner-368507.uw.r.appspot.com/";
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
        props: { date: json.date, kyoki: json.kyoki || [] },
        revalidate: 60 * 10 * 10, // 10min
    };
};
