# oVirt-Containers - master snapshot
The repository includes images descriptions and manifests definitions of ovirt
components for openshift deployment.

## MUST HAVE
Must use oc tool version 1.5.0 - https://github.com/openshift/origin/releases

WARNING: origin-clients rpm installation adds to /bin oc binary that might
be older - verify that you work with 1.5 by "oc version"

## Details
The orchestration includes engine deploymentconfig and kube-vdsm deamonset.
* For building images run "/bin/sh automation/build-artifacts.sh"
* To load deployment to openshift run "oc create -f os-manifests -R"
(for setting openshift cluster follow
https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md#linux
or https://github.com/minishift/minishift/blob/master/README.md#installation
to set up a testing instance on minishift).

## Orchestration using Minishift
- Install minishift - https://github.com/minishift/minishift
- Run the following

```
export OCTAG=v1.5.0-rc.0
export PROJECT=ovirt
export LATEST_MINISHIFT_CENTOS_ISO_BASE=$(curl -I https://github.com/minishift/minishift-centos-iso/releases/latest | grep "Location" | cut -d: -f2- | tr -d '\r' | xargs)
export MINISHIFT_CENTOS_ISO=${LATEST_MINISHIFT_CENTOS_ISO_BASE/tag/download}/minishift-centos7.iso

minishift start --memory 4096 --iso-url=$MINISHIFT_CENTOS_ISO --openshift-version=$OCTAG
export PATH=$PATH:~/.minishift/cache/oc/$OCTAG
```

### Login to openshift as system admin
```
oc login -u system:admin
```

### Create oVirt project
```
oc new-project $PROJECT --description="oVirt" --display-name="oVirt"
```

### Add administrator permissions to 'developer' user account
```
oc adm policy add-role-to-user admin developer -n $PROJECT
```

### Force a permissive security context constraints
Allows the usage of root account inside engine pod
```
oc create serviceaccount useroot
oc adm policy add-scc-to-user anyuid -z useroot
```

### Allows host IPC inside node pod
```
oc create serviceaccount privilegeduser
oc adm policy add-scc-to-user privileged -z privilegeduser
```

### Create engine and node deployments and add them to the project
Please note that the engine deployment is configured as paused
```
oc create -f os-manifests -R
```

#### To deploy just engine
```
oc create -f os-manifests/engine -R
```

#### To deploy just node
```
oc create -f os-manifests/node -R
```

### Change the hostname for the ovirt-engine deployment
According to the hostname that was assigned to the associated route
```
oc set env dc/ovirt-engine -c ovirt-engine OVIRT_FQDN=$(oc describe routes ovirt-engine | grep "Requested Host:" | cut -d: -f2 | xargs)
```

### Unpause ovirt-engine deployment
```
oc patch dc/ovirt-engine --patch '{"spec":{"paused": false}}'
```

### Provide login info
```
echo "Now you can login as developer user to the $PROJECT project, the server is accessible via web console at $(minishift console --url)"
```