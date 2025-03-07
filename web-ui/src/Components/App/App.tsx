import * as React from "react";
import { FluentProvider, webDarkTheme } from "@fluentui/react-components";
import { useStaticStyles, useStyles } from "./App.styles";
import { useStockData } from "../hooks/useStockData";
import { Chart } from "../Chart/Chart";
import { BuyIndicator } from "../BuyIndicator/BuyIndicator";
import { News } from "../News/News";
import { Filters, IFilters } from "../Filters/Filters";
import { useFilteredData } from "../hooks/useFIlteredData";
import { useStockList } from "../hooks/useStockList";

export function App() {
  return (
    <FluentProvider theme={webDarkTheme}>
      <InnerApp />
    </FluentProvider>
  );
}

function InnerApp() {
  const styles = useStyles();
  const [selectedItem, setSelectedItem] = React.useState<string | null>(null);
  const [filters, setFilters] = React.useState<IFilters | null>(null);
  useStaticStyles();

  const stockList = useStockList();
  const stockData = useStockData(stockList);
  const filteredData = useFilteredData(stockData, filters);

  const onItemSelected = React.useCallback((guid: string) => {
    setSelectedItem(guid);
  }, []);
  const onItemCleared = React.useCallback(() => {
    setSelectedItem(null);
  }, []);

  if (!stockData || stockData.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <div className={styles.app}>
      <BuyIndicator data={stockData} />
      <Filters onFilterChange={setFilters} data={stockData} />
      <Chart data={filteredData} onItemSelected={onItemSelected} />
      {selectedItem && (
        <News
          news={filteredData.find((d) => d.guid === selectedItem)}
          onClear={onItemCleared}
        />
      )}
    </div>
  );
}
