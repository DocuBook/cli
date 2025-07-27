import chalk from "chalk";

// Logging helper with styles
const log = {
  info: (msg) => console.log(chalk.cyan("! " + msg)),
  success: (msg) => console.log(chalk.green("✔ " + msg)),
  warn: (msg) => console.log(chalk.yellow("⚠️  " + msg)),
  error: (msg) => console.log(chalk.red("x " + msg)),
};

export default log;
