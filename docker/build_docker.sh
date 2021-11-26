docker image build -f Dockerfile -t forum .
if [ $? -ne 0 ]; then
    echo "Image build failed, try again"
    exit 1;
else
    while true; do
        read -p "Do you wish to clean up any unused objects [y/N] " yn
        case $yn in
            [Yy]* ) docker system prune && exit;;
            [Nn]* ) exit;;
            * ) exit;;
        esac
    done
fi
