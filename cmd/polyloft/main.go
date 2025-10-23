package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ArubikU/polyloft/internal/auth"
	"github.com/ArubikU/polyloft/internal/builder"
	"github.com/ArubikU/polyloft/internal/config"
	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/installer"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/mappings"
	"github.com/ArubikU/polyloft/internal/parser"
	"github.com/ArubikU/polyloft/internal/publisher"
	"github.com/ArubikU/polyloft/internal/repl"
	"github.com/ArubikU/polyloft/internal/searcher"
	"github.com/ArubikU/polyloft/internal/version"
)

// main provides a simple, extensible CLI entrypoint for the Polyloft project.
// Subcommands:
//   - repl: start an interactive REPL
//   - run:  run a .pf source file (placeholder pipeline)
//   - build: compile a .pf source file to a target (placeholder)
//   - version: print version information
//
// All heavy lifting is delegated to internal packages so the CLI stays thin.
func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	sub := os.Args[1]
	switch sub {
	case "repl":
		replCmd := flag.NewFlagSet("repl", flag.ExitOnError)
		prompt := replCmd.String("prompt", ">>> ", "REPL prompt string")
		_ = replCmd.Parse(os.Args[2:])
		repl.Start(os.Stdin, os.Stdout, *prompt)
	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		configFile := runCmd.String("config", "polyloft.toml", "configuration file")
		_ = runCmd.Parse(os.Args[2:])
		
		var file string
		if runCmd.NArg() < 1 {
			// No file specified, try to run the current directory as a project
			cfg, err := config.Load(*configFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "usage: polyloft run <file.pf>")
				fmt.Fprintln(os.Stderr, "Or run in a directory with polyloft.toml to use entry_point")
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}
			file = cfg.Project.EntryPoint
			if file == "" {
				fmt.Fprintln(os.Stderr, "No entry_point specified in polyloft.toml")
				os.Exit(1)
			}
		} else {
			file = runCmd.Arg(0)
		}
		
		if err := runFile(file); err != nil {
			// Use the engine's error formatter for better output
			formattedErr := engine.FormatError(err)
			fmt.Fprint(os.Stderr, formattedErr)
			os.Exit(1)
		}
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		out := buildCmd.String("o", "", "output artifact (defaults to project name)")
		configFile := buildCmd.String("config", "polyloft.toml", "configuration file")
		_ = buildCmd.Parse(os.Args[2:])

		// Load configuration
		cfg, err := config.Load(*configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Tip: Create a polyloft.toml file or use -config flag")
			os.Exit(1)
		}

		// Determine default output name if not provided
		if *out == "" {
			*out = defaultOutputName(cfg)
		} else if runtime.GOOS == "windows" {
			// Ensure Windows binaries have a runnable extension when none provided
			if filepath.Ext(*out) == "" {
				*out += ".pfx"
			}
		}

		// Build the project
		bldr := builder.New(cfg, *out)
		if err := bldr.Build(); err != nil {
			fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
			os.Exit(1)
		}
	case "install":
		installCmd := flag.NewFlagSet("install", flag.ExitOnError)
		configFile := installCmd.String("config", "polyloft.toml", "configuration file")
		globalMode := installCmd.Bool("g", false, "install packages globally")
		_ = installCmd.Parse(os.Args[2:])

		// Check if specific packages are provided as arguments
		if installCmd.NArg() > 0 {
			// Install specific packages from command line
			packages := installCmd.Args()
			
			// Try to load config for context, but don't fail if it doesn't exist
			cfg, err := config.Load(*configFile)
			if err != nil {
				// If config doesn't exist, create a minimal one for installation
				cfg = &config.Config{
					Project: config.ProjectConfig{
						Name:       "temp-install",
						Version:    "0.1.0",
						EntryPoint: "src/main.pf",
					},
				}
			}
			
			inst := installer.New(cfg)
			inst.SetGlobalMode(*globalMode)
			if err := inst.InstallPackages(packages); err != nil {
				fmt.Fprintf(os.Stderr, "Install failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Install dependencies from config file
			cfg, err := config.Load(*configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				fmt.Fprintln(os.Stderr, "Tip: Create a polyloft.toml file or use -config flag")
				os.Exit(1)
			}

			// Install dependencies
			inst := installer.New(cfg)
			inst.SetGlobalMode(*globalMode)
			if err := inst.Install(); err != nil {
				fmt.Fprintf(os.Stderr, "Install failed: %v\n", err)
				os.Exit(1)
			}
		}
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		_ = initCmd.Parse(os.Args[2:])

		// Check if polyloft.toml already exists
		if _, err := os.Stat("polyloft.toml"); err == nil {
			fmt.Fprintln(os.Stderr, "polyloft.toml already exists")
			os.Exit(1)
		}

		// Create example config
		if err := os.WriteFile("polyloft.toml", []byte(config.Example()), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create polyloft.toml: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Created polyloft.toml")
		fmt.Println("Edit the file to configure your project, then run:")
		fmt.Println("  polyloft install  - to install dependencies")
		fmt.Println("  polyloft build    - to build your project")
	case "register":
		registerCmd := flag.NewFlagSet("register", flag.ExitOnError)
		_ = registerCmd.Parse(os.Args[2:])
		
		// Interactive registration
		reader := bufio.NewReader(os.Stdin)
		
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		
		if username == "" {
			fmt.Fprintln(os.Stderr, "Username cannot be empty")
			os.Exit(1)
		}
		
		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)
		
		if email == "" {
			fmt.Fprintln(os.Stderr, "Email cannot be empty")
			os.Exit(1)
		}
		
		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		
		if password == "" {
			fmt.Fprintln(os.Stderr, "Password cannot be empty")
			os.Exit(1)
		}
		
		// Register with server
		registryURL := auth.GetRegistryURL()
		registerData := map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		}
		
		jsonData, err := json.Marshal(registerData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to prepare registration data: %v\n", err)
			os.Exit(1)
		}
		
		resp, err := http.Post(
			fmt.Sprintf("%s/api/auth/register", registryURL),
			"application/json",
			bytes.NewReader(jsonData),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to register: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Registration failed: %s\n", string(body))
			os.Exit(1)
		}
		
		fmt.Println("✓ Registration successful!")
		fmt.Println("You can now login with: polyloft login")
		
	case "login":
		loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
		_ = loginCmd.Parse(os.Args[2:])
		
		// Interactive login
		reader := bufio.NewReader(os.Stdin)
		
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		
		if username == "" {
			fmt.Fprintln(os.Stderr, "Username cannot be empty")
			os.Exit(1)
		}
		
		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		
		if password == "" {
			fmt.Fprintln(os.Stderr, "Password cannot be empty")
			os.Exit(1)
		}
		
		// Login to server
		registryURL := auth.GetRegistryURL()
		loginData := map[string]string{
			"username": username,
			"password": password,
		}
		
		jsonData, err := json.Marshal(loginData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to prepare login data: %v\n", err)
			os.Exit(1)
		}
		
		resp, err := http.Post(
			fmt.Sprintf("%s/api/auth/login", registryURL),
			"application/json",
			bytes.NewReader(jsonData),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to login: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Login failed: %s\n", string(body))
			os.Exit(1)
		}
		
		var loginResp struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse login response: %v\n", err)
			os.Exit(1)
		}
		
		// Save credentials
		creds := &auth.Credentials{
			Username: username,
			Token:    loginResp.Token,
		}
		
		if err := auth.SaveCredentials(creds); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save credentials: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✓ Successfully authenticated")
		fmt.Println("You can now use 'polyloft publish' to publish packages")
		
	case "logout":
		logoutCmd := flag.NewFlagSet("logout", flag.ExitOnError)
		_ = logoutCmd.Parse(os.Args[2:])
		
		if err := auth.ClearCredentials(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to logout: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✓ Successfully logged out")
		
	case "publish":
		publishCmd := flag.NewFlagSet("publish", flag.ExitOnError)
		configFile := publishCmd.String("config", "polyloft.toml", "configuration file")
		_ = publishCmd.Parse(os.Args[2:])
		
		// Check authentication
		if !auth.IsAuthenticated() {
			fmt.Fprintln(os.Stderr, "Not authenticated. Please run 'polyloft login' first")
			os.Exit(1)
		}
		
		// Load configuration
		cfg, err := config.Load(*configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Tip: Create a polyloft.toml file with project information")
			os.Exit(1)
		}
		
		// Publish package
		pub := publisher.New(cfg)
		if err := pub.Publish(); err != nil {
			fmt.Fprintf(os.Stderr, "Publish failed: %v\n", err)
			os.Exit(1)
		}
		
	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		_ = searchCmd.Parse(os.Args[2:])
		
		if searchCmd.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "usage: polyloft search <query>")
			os.Exit(1)
		}
		
		query := searchCmd.Arg(0)
		s := searcher.New()
		results, err := s.Search(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Search failed: %v\n", err)
			os.Exit(1)
		}
		
		if len(results) == 0 {
			fmt.Println("No packages found matching your query.")
		} else {
			fmt.Printf("Found %d package(s):\n\n", len(results))
			for _, pkg := range results {
				fmt.Printf("  %s@%s (v%s)\n", pkg.Name, pkg.Author, pkg.Version)
				if pkg.Description != "" {
					fmt.Printf("    %s\n", pkg.Description)
				}
				fmt.Println()
			}
			fmt.Println("Install with: polyloft install <package>@<author>")
		}
		
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		_ = updateCmd.Parse(os.Args[2:])
		
		// Detect platform and run appropriate update script
		scriptURL := "https://raw.githubusercontent.com/ArubikU/polyloft/main/scripts/"
		
		if runtime.GOOS == "windows" {
			// Windows PowerShell update script
			scriptURL += "update.ps1"
			fmt.Println("Downloading and running update script for Windows...")
			
			// Download and execute the PowerShell script
			resp, err := http.Get(scriptURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading update script: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			
			scriptContent, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading update script: %v\n", err)
				os.Exit(1)
			}
			
			// Save script to temp file
			tmpFile, err := os.CreateTemp("", "polyloft-update-*.ps1")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
				os.Exit(1)
			}
			defer os.Remove(tmpFile.Name())
			
			if _, err := tmpFile.Write(scriptContent); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing update script: %v\n", err)
				os.Exit(1)
			}
			tmpFile.Close()
			
			// Execute PowerShell script
			cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tmpFile.Name())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running update script: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Linux/macOS bash update script
			scriptURL += "update.sh"
			fmt.Println("Downloading and running update script for Linux/macOS...")
			
			// Download and execute the bash script
			resp, err := http.Get(scriptURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading update script: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			
			scriptContent, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading update script: %v\n", err)
				os.Exit(1)
			}
			
			// Execute bash script directly via bash -c
			cmd := exec.Command("bash", "-c", string(scriptContent))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running update script: %v\n", err)
				os.Exit(1)
			}
		}
		
	case "generate-mappings":
		genMappingsCmd := flag.NewFlagSet("generate-mappings", flag.ExitOnError)
		out := genMappingsCmd.String("o", "mappings.json", "output file path")
		root := genMappingsCmd.String("root", ".", "root directory of the project")
		_ = genMappingsCmd.Parse(os.Args[2:])
		
		fmt.Printf("Generating mappings from %s...\n", *root)
		
		gen := mappings.NewGenerator(*root)
		if err := gen.Generate(*out); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating mappings: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✓ Mappings generated successfully: %s\n", *out)
		
	case "version":
		fmt.Println(version.String())
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n\n", sub)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Polyloft CLI")
	fmt.Println()
	fmt.Println("Usage: polyloft <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  repl                  Start an interactive REPL")
	fmt.Println("  run [file.pf]         Run a Polyloft source file, or current project if no file specified")
	fmt.Println("  init                  Initialize a new project with polyloft.toml")
	fmt.Println("  build                 Build a Polyloft project to executable (requires polyloft.toml)")
	fmt.Println("  install [package]     Install project dependencies (requires polyloft.toml), or install specific package(s). Use -g for global installation")
	fmt.Println("  search <query>        Search for packages in the registry")
	fmt.Println("  register              Register a new account on the package registry")
	fmt.Println("  login                 Authenticate with the package registry")
	fmt.Println("  logout                Clear authentication credentials")
	fmt.Println("  publish               Publish package to registry (requires polyloft.toml and authentication)")
	fmt.Println("  generate-mappings     Generate mappings.json for IDE/editor support")
	fmt.Println("  update                Update Polyloft to the latest version")
	fmt.Println("  version               Print version information")
}

// runFile is a placeholder execution pipeline that shows where
// lexing/parsing/execution will be wired in the future.
func runFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	source := string(b)

	// Tokenize
	lx := &lexer.Lexer{}
	items := lx.Scan(b)

	// Parse (with filename and source for better errors)
	p := parser.NewWithSource(items, path, source)
	prog, err := p.Parse()
	if err != nil {
		return err
	}

	// Eval with file context and source for better error messages
	packageName := filepath.Dir(path)
	_, err = engine.EvalWithContextAndSource(prog, engine.Options{Stdout: os.Stdout}, path, packageName, source)
	return err
}

// defaultOutputName builds a sensible default artifact name based on config and OS.
func defaultOutputName(cfg *config.Config) string {
	name := cfg.Project.Name
	if name == "" {
		name = "polyloft-app"
	}
	return name + ".pfx"
}
