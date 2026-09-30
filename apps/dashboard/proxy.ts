import { NextResponse, type NextRequest } from "next/server";

// Cheap presence check; the API validates the session on every request.
export function proxy(request: NextRequest) {
  if (request.cookies.has("session")) {
    return NextResponse.next();
  }
  return NextResponse.redirect(new URL("/login", webOrigin(request)));
}

// The dashboard is usually reached through the web app's rewrite, so send
// users back to the web app's origin rather than the dashboard deployment.
function webOrigin(request: NextRequest) {
  if (process.env.WEB_URL) return process.env.WEB_URL;
  const host = request.headers.get("x-forwarded-host");
  const proto = request.headers.get("x-forwarded-proto") ?? "https";
  return host ? `${proto}://${host}` : request.nextUrl.origin;
}

export const config = {
  // "/" is listed separately so the basePath root (/dashboard) is matched too.
  matcher: ["/", "/((?!_next/static|_next/image|favicon.ico).*)"],
};
