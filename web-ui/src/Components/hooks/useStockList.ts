import React from "react";
import { SERVER_BASE_URL } from "../../constants";

const dataUrl = `${SERVER_BASE_URL}subjects`;

export function useStockList(): string[] {
  const [stockList, setStockList] = React.useState<string[]>([]);
  React.useEffect(() => {
    fetch(dataUrl)
      .then((d) => d.json())
      .then(setStockList);
  }, []);

  return stockList;
}
