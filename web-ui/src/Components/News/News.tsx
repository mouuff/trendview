import React from "react";
import { IStockData } from "../hooks/useStockData";
import { useStyles } from "./News.styles";
import { Dismiss20Regular } from "@fluentui/react-icons";

export function News({
  news,
  onClear,
}: {
  news: IStockData;
  onClear: () => void;
}) {
  const styles = useStyles();

  return (
    <div className={styles.news}>
      <button className={styles.closeBtn} onClick={onClear}>
        <Dismiss20Regular />
      </button>
      <h3 className={styles.content}>{news.title}</h3>
      <span className={styles.date}>
        {news.source}
        {" - "}
        {new Date(news.dateTime).toDateString()}
      </span>
    </div>
  );
}
