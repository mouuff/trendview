import { makeStyles, tokens } from "@fluentui/react-components";

export const useStyles = makeStyles({
  news: {
    display: "flex",
    position: "absolute",
    left: "50%",
    top: "50%",
    transform: "translate(-50%, -50%)",
    flexDirection: "column",
    padding: "5px 10px 10px 10px",
    backgroundColor: tokens.colorNeutralBackground2,
    color: tokens.colorNeutralForeground2,
    borderRadius: "5px",
    border: `1px solid ${tokens.colorNeutralStroke1}`,
    width: "50%",
    height: "auto",
  },
  closeBtn: {
    position: "absolute",
    top: "5px",
    right: "5px",
    backgroundColor: tokens.colorNeutralBackground2,
    color: tokens.colorNeutralForeground2,
    padding: "8px",
    marginBottom: "10px",
    border: `1px solid transparent`,
    boxSizing: "border-box",
    ":hover": {
      cursor: "pointer",
      border: `1px solid ${tokens.colorNeutralStroke1}`,
      borderRadius: "4px",
    },
  },
  content: {
    fontSize: "20px",
  },
  date: {
    display: "flex",
    justifyContent: "flex-end",
    fontSize: "14px",
  },
});
