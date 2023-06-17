export default defineNuxtPlugin({
  name: "insert_sfc_code_from_server",
  enforce: "default",
  hooks: {
    "page:finish": async (page) => {
      const test = {
        mountId: "some-id",
        sfcCode:
          "<template><p>{{ hello }}</p></template><script>export default { data() { return { hello: 'hello' } } };</script><style scoped>& p { color: red; }</style>",
      };
      const response = await fetch("/api", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(test),
      });
      const data = await response.json();

      const script = document.createElement("script");
      script.id = `script-code-${test.mountId}`;
      script.innerHTML = data.appCode;
      script.type = "module";
      script.onerror = () => {
        console.log("Error occurred while loading script");
      };
      document.body.appendChild(script);

      const style = document.createElement("style");
      style.id = `style-code-${test.mountId}`;
      style.textContent = data.styleCode.reduce(
        (styleContent: string, rule: string) => {
          styleContent += `${rule}\n`;
          return styleContent;
        },
        ""
      );
      document.head.appendChild(style);
    },
  },
});
