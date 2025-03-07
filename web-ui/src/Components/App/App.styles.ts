import { makeStaticStyles, makeStyles } from "@fluentui/react-components";

// Some static styles for the top level elements
export const useStaticStyles = makeStaticStyles({
  body: {
    margin: 0,
    overflow: "hidden"
  },
});

export const useStyles = makeStyles({
  app: {
    display: "flex",
    flexDirection: "column",
    // set app size
    height: "100vh",
    // use a dark background color
    backgroundColor: "#111",
    // default font family
    fontFamily: "Arial, sans-serif",
    fontSize: "35px",
    padding: "0px 20px",
    overflow: "hidden"
  },
});
