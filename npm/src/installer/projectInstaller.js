import path from "path";
import fs from "fs";
import ora from "ora";
import chalk from "chalk";
import { execSync } from "child_process";
import { fileURLToPath } from "url";
import log from "../utils/logger.js";
import { configurePackageManager } from "../utils/packageManager.js";
import { displayManualSteps, simulateInstallation, displayNextSteps } from "../utils/display.js";

/**
 * Creates a new DocuBook project.
 * @param {Object} options - Installation options.
 */
export async function createProject({ directoryName, packageManager, version, installNow }) {
  const projectPath = path.resolve(process.cwd(), directoryName);

  if (fs.existsSync(projectPath)) {
    throw new Error(`Directory "${directoryName}" already exists.`);
  }

  log.info(`Creating a new DocuBook project in ${chalk.cyan(projectPath)}...`);

  const spinner = ora("Setting up project files...").start();

  try {
    // 1. Create project directory and copy template files
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = path.dirname(__filename);
    const templatePath = path.join(__dirname, "../dist");
    copyDirectoryRecursive(templatePath, projectPath);
    spinner.succeed("Project files created.");

    // 2. Configure package manager specific settings
    spinner.start("Configuring package manager...");
    configurePackageManager(packageManager, projectPath);
    spinner.succeed("Package manager configured.");

    // 3. Update package.json
    spinner.start("Updating package.json...");
    const pkgPath = path.join(projectPath, "package.json");
    if (fs.existsSync(pkgPath)) {
      const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf-8"));
      pkg.name = directoryName; // Set project name
      pkg.packageManager = `${packageManager}@${version}`;
      fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
    }
    spinner.succeed("package.json updated.");

    log.success(`DocuBook project ready to go!`);

    if (installNow) {
      await installDependencies(directoryName, packageManager, projectPath);
      await simulateInstallation();
      displayNextSteps(directoryName, packageManager);
    } else {
      displayManualSteps(directoryName, packageManager);
    }
  } catch (err) {
    spinner.fail("Failed to create project.");
    // Cleanup created directory on failure
    if (fs.existsSync(projectPath)) {
      fs.rmSync(projectPath, { recursive: true, force: true });
    }
    throw err;
  }
}

/**
 * Recursively copies a directory.
 * @param {string} source - Source directory path.
 * @param {string} destination - Destination directory path.
 */
function copyDirectoryRecursive(source, destination) {
  if (!fs.existsSync(destination)) {
    fs.mkdirSync(destination, { recursive: true });
  }

  const entries = fs.readdirSync(source, { withFileTypes: true });
  for (const entry of entries) {
    const srcPath = path.join(source, entry.name);
    const destPath = path.join(destination, entry.name);

    if (entry.isDirectory()) {
      copyDirectoryRecursive(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

/**
 * Installs project dependencies.
 * @param {string} directoryName - Project directory name.
 * @param {string} packageManager - Package manager to use.
 * @param {string} projectPath - Path to the project directory.
 */
async function installDependencies(directoryName, packageManager, projectPath) {
  log.info("Installing dependencies...");
  const installSpinner = ora(`Running ${chalk.cyan(`${packageManager} install`)}...`).start();

  try {
    execSync(`${packageManager} install`, { cwd: projectPath, stdio: "ignore" });
    installSpinner.succeed("Dependencies installed successfully.");
  } catch (error) {
    installSpinner.fail("Failed to install dependencies.");
    displayManualSteps(directoryName, packageManager);
    throw new Error("Dependency installation failed.");
  }
}