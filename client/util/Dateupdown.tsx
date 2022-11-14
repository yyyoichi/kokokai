import { Box } from "@mui/material";
import Stack from "@mui/material/Stack";
import { useRouter } from "next/router";
import { getYYYYMMDD } from "../src/FormatDate";

export default function Dateupdown({ dateString }: { dateString: string }) {
    const date = new Date(dateString);
    const router = useRouter();
    const getPath = (di: Date) => {
        return `/daykyoki/${getYYYYMMDD(di)}`;
    };
    console.log(date);
    const up = (
        <div
            onClick={() => {
                const upd = new Date(date);
                upd.setDate(upd.getDate() + 1);
                router.push(getPath(upd));
            }}
        >
            {">"}
        </div>
    );
    const down = (
        <div
            onClick={() => {
                const upd = new Date(date);
                upd.setDate(upd.getDate() - 1);
                router.push(getPath(upd));
            }}
        >
            {"<"}
        </div>
    );
    return (
        <Stack direction="row" spacing={2} justifyContent={"center"} my={2}>
            {down}
            <Box>{dateString}</Box>
            {up}
        </Stack>
    );
}
