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
    .argument("[directory]", "The name of the project directory")
    .action(async (directory) => {
      try {
        displayIntro();
        const userInput = await collectUserInput(directory);

        // Pass all user input AND the DocuBook version to createProject
        await createProject({
          ...userInput,
          docubookVersion: version, // Add DocuBook version here
        });
      } catch (err) {
        log.error(err.message || "An unexpected error occurred.");
        process.exit(1);
      }
    });

  return program;
}
