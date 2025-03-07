import * as path from "path";
import HtmlWebpackPlugin from "html-webpack-plugin";
import { Configuration, DefinePlugin } from "webpack";
import MiniCssExtractPlugin from "mini-css-extract-plugin";
import { GriffelCSSExtractionPlugin } from "@griffel/webpack-extraction-plugin";

type Args = {
  mode: "development" | "production";
  env: Record<string, unknown>;
};

export default (_: never, { mode }: Args): Configuration => {
  return {
    entry: "./src/index.tsx",
    mode,
    resolve: {
      extensions: [".tsx", ".ts", ".js"],
    },
    output: {
      filename: "[name].js",
      path: path.resolve(__dirname, "dist"),
    },
    module: {
      rules: [
        // Add the TypeScript loader
        {
          test: /\.tsx?$/,
          use: "ts-loader",
          exclude: /node_modules/,
        },
        // add griffel css extraction loader
        {
          test: /\.(js|ts|tsx)$/,
          include: [
            path.resolve(__dirname, "src/components"),
            /\/node_modules\/@fluentui\//,
          ],
          use: {
            loader: GriffelCSSExtractionPlugin.loader,
          },
        },
        // Add the griffel loader (for CSS in JS)
        {
          test: /\.styles\.(ts|tsx)$/,
          exclude: /node_modules/,
          use: {
            loader: "@griffel/webpack-loader",
            options: {
              babelOptions: {
                presets: ["@babel/preset-typescript"],
              },
            },
          },
        },
        // "css-loader" is required to handle produced CSS assets by Griffel
        // "mini-css-extract-plugin" is required to extract CSS into separate files
        {
          test: /\.css$/,
          use: [MiniCssExtractPlugin.loader, "css-loader"],
        },
      ],
    },
    plugins: [
      new DefinePlugin({
        __DEV__: mode === "development",
      }),
      new HtmlWebpackPlugin({
        template: "public/index.html",
      }),
      new MiniCssExtractPlugin({
        ignoreOrder: true,
      }),
      new GriffelCSSExtractionPlugin(),
    ],
  };
};
