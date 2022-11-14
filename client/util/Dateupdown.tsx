import { Box } from "@mui/material";
import Stack from "@mui/material/Stack";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import { useRouter } from "next/router";
import { getYYYYMMDD } from "../src/FormatDate";

export default function Dateupdown({ dateString }: { dateString: string }) {
    const date = new Date(dateString);
    const router = useRouter();
    const getPath = (di: Date) => {
        return `/daykyoki/${getYYYYMMDD(di)}`;
    };
    const up = (
        <div
            onClick={() => {
                const upd = new Date(date);
                upd.setDate(upd.getDate() + 1);
                router.push(getPath(upd));
            }}
        >
            <ArrowForwardIosIcon />
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
            <ArrowBackIosNewIcon />
        </div>
    );
    return (
        <Stack direction="row" spacing={2} justifyContent={"center"} my={2}>
            {down}
            <Box fontWeight={500} fontSize={18}>
                {dateString}
            </Box>
            {up}
        </Stack>
    );
}
