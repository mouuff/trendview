import React from "react";
import { RELATED_RATING_CUTOFF, SERVER_BASE_URL } from "../../constants";

// check if dev

type IBtcDataResponse = Record<
  string,
  {
    Title: string;
    Content: string;
    DateTime: string;
    Link: string;
    GUID: string;
    Source: string;
    Results: {
      Confidence: number;
      Relevance: number;
    };
  }
>;

export interface IStockData {
  title: string;
  content: string;
  dateTime: string;
  link: string;
  guid: string;
  source: string;
  rating: number;
  stockName: string;
}

const fetchDataForStock = (stock: string): Promise<IStockData[]> =>
  fetch(`${SERVER_BASE_URL}itemsBySubject?subject=${stock}`)
    .then((d) => d.json())
    .then(({ Items: data }: { Items: IBtcDataResponse }) =>
      Object.keys(data)
        .filter((key) => data[key].Results.Relevance >= RELATED_RATING_CUTOFF)
        .map((key) => ({
          title: data[key].Title,
          content: data[key].Content,
          dateTime: data[key].DateTime,
          link: data[key].Link,
          guid: data[key].GUID,
          source: data[key].Source,
          rating: data[key].Results.Confidence,
          stockName: stock,
        }))
    );

export function useStockData(stockList: string[]) {
  const [stockData, setStockData] = React.useState<IStockData[]>([]);

  // on mount, fetch
  React.useEffect(() => {
    Promise.all<IStockData[]>(stockList.map((s) => fetchDataForStock(s)))
      .then((datas) => datas.reduce((acc, d) => acc.concat(d), []))
      .then(setStockData);
  }, [stockList]);

  return stockData;
}
