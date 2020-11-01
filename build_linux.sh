# make built folder
echo "Make built folder"
mkdir built
built_folder=built/`cat version`_linux_`date +%Y%m%d_%H%M%S`
mkdir $built_folder

# go build
echo "Run go build"
go build -o "$built_folder/blog.run" .

# copy files to built folder
echo "Copy help.txt"
cp help.txt $built_folder/
echo "Copy drawer/web/"
mkdir $built_folder/drawer
cp -r drawer/web $built_folder/drawer/

# make sqlite database. NEED SQLITE3!!!
echo "Make database"
sqlite3 $built_folder/blog.db ".read whole_database.sql"
echo "Build done. Check $built_folder."
