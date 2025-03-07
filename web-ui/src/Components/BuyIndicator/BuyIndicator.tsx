import React from "react";
import { IStockData } from "../hooks/useStockData";
import { BUY_THRESHOLD, DONT_BUY_THRESHOLD } from "../../constants";

// return only items which are dated from today
const filterDaily = (data: IStockData[]) => {
  const today = new Date().toISOString().split("T")[0];
  return data.filter((d) => d.dateTime.includes(today));
};

const buygifs = [
  "https://media.tenor.com/dblb_XKGVC4AAAAj/pepe-the-frog-sad.gif",
  "https://media.tenor.com/dblb_XKGVC4AAAAM/pepe-the-frog-sad.gif",
  "https://media.tenor.com/HcSqaQeyukYAAAAj/pepeblankies-pepenet.gif",
  "https://media.tenor.com/bw4mKWWIA5IAAAAj/peepo-pepe-the-frog.gif",
  "https://media.tenor.com/6QXSPtLY0cMAAAAi/pepe-hug.gif",
  "https://media.tenor.com/Xh5lvnAh1IsAAAAj/pepe-love-pepe.gif",
  "https://media.tenor.com/M6LfbrDeUEkAAAAj/yahh.gif",
  "https://media.tenor.com/9_RPqsHQN6cAAAAj/ppp.gif",
  "https://media.tenor.com/yFz1Phe5SlwAAAAj/gaming-pepe.gif",
  "https://media.tenor.com/6wl9HBQjGxIAAAAj/pepe.gif",
  "https://media.tenor.com/zn2-vlKrzWUAAAAj/pepe-pain.gif",
  "https://media.tenor.com/VR-bsaKAAXEAAAAj/pepe-cute.gif",
  "https://media.tenor.com/wT-64avQxegAAAAj/i-love-you-pepe.gif",
  "https://media.tenor.com/S1QLV3bVI6oAAAAj/lubbers-pepe.gif",
  "https://media.tenor.com/jCqnf9SBcVsAAAAj/peepo-sit-sitti.gif",
  "https://media.tenor.com/3OVWEW4Cw50AAAAj/feelsdankman-feelsdank.gif",
];

const getRandomBuyGif = () => {
  return buygifs[Math.floor(Math.random() * buygifs.length)];
};

export function BuyIndicator({ data }: { data: IStockData[] }) {
  const rating = React.useMemo(() => {
    if (data.length === 0) {
      return -1;
    }
    return Math.round(data.reduce((acc, d) => acc + d.rating, 0) / data.length);
  }, [data]);

  const indicatorContent = React.useMemo(() => {
    const shouldBy =
      rating > BUY_THRESHOLD
        ? "Yes"
        : rating < DONT_BUY_THRESHOLD
        ? "No"
        : "Unsure";

    if (shouldBy === "Unsure") {
      return <span>"Unsure, ᕦ(ò_óˇ)ᕤ"</span>;
    }

    if (shouldBy === "Yes") {
      return (
        <span>
          BUY
          <img src={getRandomBuyGif()} />
        </span>
      );
    }

    return `Don't buy`;
  }, [rating]);

  return (
    <div>
      <h2>Buy Indicator</h2>
      <p>
        {indicatorContent}
        (rating: {rating}/100)
      </p>
    </div>
  );
}
