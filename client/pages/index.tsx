import { Box, Button, Container, Stack } from "@mui/material";
import Head from "next/head";
import { env } from "process";
import { useEffect, useState } from "react";
import { Anton } from "@next/font/google";
import Link from "next/link";

const anton = Anton({ weight: "400" });

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
            <Box
                sx={{
                    background:
                        "linear-gradient(to right bottom, #FFF, #8BC6EC, #9599E2)",
                }}
            >
                <Container
                    maxWidth="sm"
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        minHeight: "100vh",
                        alignItems: "center",
                    }}
                >
                    <Box py={4} pl={2}>
                        <h1 className={anton.className}>Collokai</h1>
                    </Box>
                    <Box py={10}>
                        <p>国会議事録APIからコロケーションで</p>
                        <p>その日のトピックをざっくり把握する</p>
                    </Box>
                    <Box mt={"auto"} pb={10}>
                        <Button variant="contained">
                            <Link href={"/daykyoki/20221102"}>
                                トピックをみる
                            </Link>
                        </Button>
                    </Box>
                </Container>
            </Box>
        </div>
    );
}
