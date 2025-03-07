import { makeStyles, tokens } from "@fluentui/react-components";

export const useStyles = makeStyles({
    chart: {
        width: "100%",
        maxHeight: "800px",
        color: tokens.colorNeutralForeground1,
        backgroundColor: "rgba(0,0,0,0.5)",
        "& text": {
            fill: `${tokens.colorNeutralForeground1} !important`,
        }
    },
    label: {
        fontSize: "12px"
    }
})