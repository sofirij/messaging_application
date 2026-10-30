import { readFileSync } from "fs"
import { importSPKI } from "jose/key/import"

const publicPem = readFileSync("../public.pem", "utf8")
export const jwtPublicKey = await importSPKI(publicPem, "EdDSA")