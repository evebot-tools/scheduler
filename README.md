# #APP#

# Next Steps

After creating a repo from the template run the following commands at the root of the new repo

```shell
find . -type f -exec sed -i 's/#APP#/<APP_NAME>/g' {} +
find . -type f -exec sed -i 's/#REPO#/<REPO_NAME>/g' {} +

```
