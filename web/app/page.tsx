import { Suspense } from "react";
import HealthCheck from "./health-check";

// 疎通確認用。サーバー（コンテナ内）から Go API を叩く経路を試す。
// ここで使う URL は compose ネットワーク内のホスト名なので API_URL。
// localhost では api に届かない（localhost は web コンテナ自身を指す）。
async function ServerHealth() {
  const apiUrl = process.env.API_URL ?? "http://api:8080";
  let result: string;
  try {
    const res = await fetch(`${apiUrl}/healthz`, { cache: "no-store" });
    result = `${res.status} ${(await res.text()).trim()}`;
  } catch (err: unknown) {
    result = `NG: ${String(err)}`;
  }

  return (
    <dl>
      <dt>web コンテナ → api（API_URL）</dt>
      <dd>
        <code>{apiUrl}/healthz</code>
      </dd>
      <dd>
        <code>{result}</code>
      </dd>
    </dl>
  );
}

export default function Page() {
  return (
    <main style={{ padding: "2rem", fontFamily: "system-ui, sans-serif", lineHeight: 1.7 }}>
      <h1>warican</h1>
      <p>環境の疎通確認ページ。アプリを書き始めたらこのファイルごと消してよい。</p>
      <Suspense fallback={<p>checking server side...</p>}>
        <ServerHealth />
      </Suspense>
      <HealthCheck />
    </main>
  );
}
