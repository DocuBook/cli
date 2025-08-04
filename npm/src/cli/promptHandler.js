import prompts from "prompts";
import { detectDefaultPackageManager, getPackageManagerVersion } from "../utils/packageManager.js";
import log from "../utils/logger.js";
import chalk from "chalk";

/**
 * Collects user input for project creation
 * @param {string} [cliProvidedDir] - The directory name provided via CLI argument.
 * @returns {Promise<Object>} User answers
 */
export async function collectUserInput(cliProvidedDir) {
  const defaultPackageManager = detectDefaultPackageManager();
  let answers = {
    directoryName: cliProvidedDir
  };

  const questions = [
    {
      type: cliProvidedDir ? null : "text",
      name: "directoryName",
      message: "What is your project named?",
      initial: "docubook",
      validate: (name) => name.trim().length > 0 ? true : "Project name cannot be empty.",
    },
    {
      type: "select",
      name: "packageManager",
      message: "Select your package manager:",
      choices: [
        { title: "npm", value: "npm" },
        { title: "pnpm", value: "pnpm" },
        { title: "yarn", value: "yarn" },
        { title: "bun", value: "bun" },
      ],
      initial: ["npm", "pnpm", "yarn", "bun"].indexOf(defaultPackageManager),
    },
    {
      type: (prev) => (prev !== 'yarn' ? "confirm" : null),
      name: "installNow",
      message: "Would you like to install dependencies now?",
      initial: true,
    },
  ];

  const promptAnswers = await prompts(questions, {
    onCancel: () => {
      log.error("Scaffolding cancelled.");
      process.exit(0);
    },
  });

  answers = { ...answers, ...promptAnswers };

  // Validate the selected package manager
  const pmVersion = getPackageManagerVersion(answers.packageManager);
  if (!pmVersion) {
    throw new Error(`${chalk.bold(answers.packageManager)} is not installed on your system. Please install it to continue.`);
  }

  // Return all answers, including the package manager version with a clear property name
  return {
    ...answers,
    directoryName: answers.directoryName.trim(),
    pmVersion,
  };
}
