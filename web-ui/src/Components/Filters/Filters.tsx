import {
  Button,
  Field,
  Menu,
  MenuItemCheckbox,
  MenuList,
  MenuPopover,
  MenuProps,
  MenuTrigger,
} from "@fluentui/react-components";
import { DatePicker } from "@fluentui/react-datepicker-compat";
import React from "react";
import { useStyles } from "./Filters.styles";
import { IStockData } from "../hooks/useStockData";

export type IFilters = Record<
  "stock" | "source" | "date-from" | "date-to",
  string[]
>;

interface Props {
  onFilterChange: (filters: IFilters) => void;
  data: IStockData[];
}

const inputFrom = () =>
  document.getElementById("date-from") as HTMLInputElement;
const inputTo = () => document.getElementById("date-to") as HTMLInputElement;

const updateTime = (
  nbDaysAgo: number,
  onDateChange: (input: "from" | "to", date: Date) => void
) => {
  const today = new Date();
  const startDate = new Date();
  startDate.setDate(startDate.getDate() - nbDaysAgo);

  onDateChange("from", startDate);
  onDateChange("to", today);
  setTimeout(() => {
    inputFrom().value = startDate.toDateString();
    inputTo().value = today.toDateString();
  }, 1);
};

export function Filters({ data, onFilterChange }: Props) {
  const [filters, setFilters] = React.useState<IFilters>({
    stock: [],
    source: [],
    "date-from": [],
    "date-to": [],
  });
  const styles = useStyles();

  const [stocks, sources] = React.useMemo(() => {
    const stocks = Object.keys(
      data.reduce<Record<string, string>>((acc, d) => {
        acc[d.stockName] = d.stockName;
        return acc;
      }, {})
    );
    const sources = Object.keys(
      data.reduce<Record<string, string>>((acc, d) => {
        acc[d.source] = d.source;
        return acc;
      }, {})
    );

    return [stocks, sources];
  }, [data]);

  const onStockChange: MenuProps["onCheckedValueChange"] = (
    _: any,
    { name, checkedItems }: any
  ) => {
    const shouldRemove = checkedItems.length === 0;

    if (shouldRemove) {
      setFilters((f) => ({
        ...f,
        stock: f.stock.filter((s) => s !== name),
      }));
    } else {
      setFilters((f) => ({
        ...f,
        stock: f.stock.concat(name),
      }));
    }
  };

  const onSourceChange: MenuProps["onCheckedValueChange"] = (
    _: any,
    { name, checkedItems }: any
  ) => {
    const shouldRemove = checkedItems.length === 0;

    if (shouldRemove) {
      setFilters((f) => ({
        ...f,
        source: f.source.filter((s) => s !== name),
      }));
    } else {
      setFilters((f) => ({
        ...f,
        source: f.source.concat(name),
      }));
    }
  };

  const onDateChange = React.useCallback((input: "from" | "to", date: Date) => {
    setFilters((filters) => ({
      ...filters,
      [`date-${input}`]: [date.toISOString()],
    }));
  }, []);

  React.useEffect(() => {
    onFilterChange(filters);
  }, [filters]);

  return (
    <div className={styles.filters}>
      <Menu>
        <MenuTrigger disableButtonEnhancement>
          <Button className={styles.menuBtn}>Stock filter</Button>
        </MenuTrigger>

        <MenuPopover>
          <MenuList
            checkedValues={filters.stock.reduce((acc, f) => {
              acc[f] = [f];
              return acc;
            }, {} as Record<string, string[]>)}
            onCheckedValueChange={onStockChange}
          >
            {stocks.map((stock) => (
              <MenuItemCheckbox key={stock} name={stock} value={stock}>
                {stock}
              </MenuItemCheckbox>
            ))}
          </MenuList>
        </MenuPopover>
      </Menu>
      <Menu>
        <MenuTrigger disableButtonEnhancement>
          <Button className={styles.menuBtn}>Source filter</Button>
        </MenuTrigger>

        <MenuPopover>
          <MenuList
            checkedValues={filters.source.reduce((acc, f) => {
              acc[f] = [f];
              return acc;
            }, {} as Record<string, string[]>)}
            onCheckedValueChange={onSourceChange}
          >
            {sources.map((source) => (
              <MenuItemCheckbox key={source} name={source} value={source}>
                {source}
              </MenuItemCheckbox>
            ))}
          </MenuList>
        </MenuPopover>
      </Menu>
      <Field>
        <DatePicker
          id="date-from"
          placeholder="Start date..."
          onSelectDate={(date) => onDateChange("from", date)}
        />
      </Field>
      <Field>
        <DatePicker
          id="date-to"
          placeholder="End date..."
          onSelectDate={(date) => onDateChange("to", date)}
        />
      </Field>
      <Button
        onClick={() => {
          updateTime(1, onDateChange);
        }}
      >
        Last 24h
      </Button>

      <Button
        onClick={() => {
          updateTime(7, onDateChange);
        }}
      >
        Last 7 days
      </Button>
      <Button
        onClick={() => {
          updateTime(30, onDateChange);
        }}
      >
        Last 30 days
      </Button>
    </div>
  );
}
