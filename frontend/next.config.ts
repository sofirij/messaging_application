import type { NextConfig } from "next";

const nextConfig: NextConfig = {
	output: "export",
	skipTrailingSlashRedirect: true,
	reactCompiler: true,
	images: {
		loader: "custom",
		loaderFile: "./src/image/loader.ts"
	}
};

export default nextConfig;
