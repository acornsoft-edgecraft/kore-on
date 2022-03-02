echo "package conf\n\nconst Template = \`#koreon.toml" > pkg/conf/template.go
cat docs/koreon.sample.toml >> pkg/conf/template.go
echo "\n\`" >> pkg/conf/template.go
