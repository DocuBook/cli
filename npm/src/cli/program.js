import { program } from "commander";
import { displayIntro } from "../utils/display.js";
import { collectUserInput } from "./promptHandler.js";
import { createProject } from "../installer/projectInstaller.js";
import log from "../utils/logger.js";

/**
 * Initializes the CLI program
 * @param {string} version - Package version
 */
export function initializeProgram(version) {
  program
    .version(version)
    .description("CLI to create a new DocuBook project")
    // 1. Define optional argument here (handled below with .argument)
    .argument("[directory]", "The name of the project directory")
    .action(async (directory) => { // 2. Capture argument in action function
      try {
        displayIntro();
        // 3. Pass the argument to the prompt function
        const options = await collectUserInput(directory);

        await createProject(options);
      } catch (err) {
        log.error(err.message || "An unexpected error occurred.");
        process.exit(1);
      }
    });

  return program;
}