import { NextResponse, type NextRequest } from "next/server";

// Cheap presence check; the API validates the session on every request.
export function proxy(request: NextRequest) {
  if (request.cookies.has("session")) {
    return NextResponse.next();
  }
  return NextResponse.redirect(new URL("/login", webOrigin(request)));
}

// Locally the dashboard is reached through the web app's rewrite, so use the
// forwarded host to send users back to the web app rather than to :3001.
function webOrigin(request: NextRequest) {
  const host = request.headers.get("x-forwarded-host");
  const proto = request.headers.get("x-forwarded-proto") ?? "https";
  return host ? `${proto}://${host}` : request.nextUrl.origin;
}

export const config = {
  // "/" is listed separately so the basePath root (/dashboard) is matched too.
  matcher: ["/", "/((?!_next/static|_next/image|favicon.ico).*)"],
};
