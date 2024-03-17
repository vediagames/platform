version=$1
serviceName=$2

pulumiName=$serviceName"Version"

sed -i "s/\($pulumiName: \).*/$pulumiName: $version/g" gcp/pulumi/app/Pulumi.prod.yaml
sed -i "s/vediagames\/vg_$serviceName:.*/vediagames\/vg_$serviceName:$version/g" vedia-local-k8s/vediagames/base/$serviceName/deployment.yaml
sed -i "s/vediagames\/vg_$serviceName:.*/vediagames\/vg_$serviceName:$version/g" vg-k3s-01/vediagames/$serviceName/deployment.yaml
