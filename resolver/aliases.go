package resolver

// Aliases maps human-friendly language names to their exact Arch package names.
var Aliases = map[string]string{
	// Core Web / Scripting
	"node":   "nodejs",
	"python": "python",
	"python2": "python2", // Usually AUR now
	"ruby":   "ruby",
	"php":    "php",
	"lua":    "lua",

	// Systems / Compiled
	"go":     "go",
	"rust":   "rust",
	"zig":    "zig",
	"nim":    "nim",
	"c":      "gcc",
	"cpp":    "gcc",
	"clang":  "clang",

	// The Java Ecosystem (Arch maintains these specific names)
	"java":   "jre-openjdk",   // Always the bleeding edge latest
	"java8":  "jre8-openjdk",
	"java11": "jre11-openjdk",
	"java17": "jre17-openjdk",
	"java21": "jre21-openjdk",
	"java22": "jre22-openjdk",

	// Dotnet
	"dotnet": "dotnet-sdk",
	"csharp": "dotnet-sdk",
}