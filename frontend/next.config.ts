import type { NextConfig } from "next";

const nextConfig: NextConfig = {
	output: "export",
	reactCompiler: true,
	images: {
		loader: "custom",
		loaderFile: "./src/image/loader.ts"
	}
};

export default nextConfig;
