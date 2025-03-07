import { IStockData } from "../hooks/useStockData";
import React from "react";
import { axisClasses, ScatterChart } from "@mui/x-charts";
import { useStyles } from "./Chart.styles";

const chartSetting = {
  yAxis: [
    {
      label: "rainfall (mm)",
    },
  ],
  sx: {
    [`.${axisClasses.left} .${axisClasses.label}`]: {
      transform: "translate(-20px, 0)",
    },
  },
};

const stringTocolorHash = (str: string) => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  let color = "#";
  for (let i = 0; i < 3; i++) {
    let value = (hash >> (i * 8)) & 0xff;
    color += ("00" + value.toString(16)).substr(-2);
  }
  return color;
};

export function Chart({
  data,
  onItemSelected,
}: {
  data: IStockData[];
  onItemSelected: (guid: string) => void;
}) {
  const styles = useStyles();
  if (data.length === 0) {
    return <div>No data found</div>;
  }
  return (
    <ScatterChart
      className={styles.chart}
      onItemClick={(_: any, params: any) => {
        onItemSelected(params.seriesId);
      }}
      series={data.map((d) => ({
        label: d.title,
        color: stringTocolorHash(d.stockName),
        id: d.guid,
        data: [
          {
            x: new Date(d.dateTime).getTime() / 1000000,
            y: d.rating,
            id: d.guid,
          },
        ],
      }))}
      {...chartSetting}
      slotProps={{
        legend: { hidden: true },
      }}
    />
  );
}
