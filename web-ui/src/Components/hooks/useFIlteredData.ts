import { IFilters } from "../Filters/Filters";
import { IStockData } from "./useStockData";

export function useFilteredData(
  data: IStockData[],
  filters: IFilters | null
): IStockData[] {
  return data.filter((d) => {
    if (!filters) {
      return true;
    }

    if (filters.stock.length > 0 && !filters.stock.includes(d.stockName)) {
      return false;
    }

    if (filters.source.length > 0 && !filters.source.includes(d.source)) {
      return false;
    }

    if (filters["date-from"].length > 0) {
      const fromDate = new Date(filters["date-from"][0]);
      if (fromDate > new Date(d.dateTime)) {
        return false;
      }
    }
    if (filters["date-to"].length > 0) {
      const toDate = new Date(filters["date-to"][0]);
      if (toDate < new Date(d.dateTime)) {
        return false;
      }
    }
    return true;
  });
}
