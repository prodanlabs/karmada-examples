package crds

const ClusterOverridePolicies = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.6.2
  creationTimestamp: null
  name: clusteroverridepolicies.policy.karmada.io
spec:
  group: policy.karmada.io
  names:
    kind: ClusterOverridePolicy
    listKind: ClusterOverridePolicyList
    plural: clusteroverridepolicies
    shortNames:
    - cop
    singular: clusteroverridepolicy
  scope: Cluster
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        description: ClusterOverridePolicy represents the cluster-wide policy that
          overrides a group of resources to one or more clusters.
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: Spec represents the desired behavior of ClusterOverridePolicy.
            properties:
              overriders:
                description: Overriders represents the override rules that would apply
                  on resources
                properties:
                  argsOverrider:
                    description: ArgsOverrider represents the rules dedicated to handling
                      container args
                    items:
                      description: CommandArgsOverrider represents the rules dedicated
                        to handling command/args overrides.
                      properties:
                        containerName:
                          description: The name of container
                          type: string
                        operator:
                          description: Operator represents the operator which will
                            apply on the command/args.
                          enum:
                          - add
                          - remove
                          type: string
                        value:
                          description: Value to be applied to command/args. Items
                            in Value which will be appended after command/args when
                            Operator is 'add'. Items in Value which match in command/args
                            will be deleted when Operator is 'remove'. If Value is
                            empty, then the command/args will remain the same.
                          items:
                            type: string
                          type: array
                      required:
                      - containerName
                      - operator
                      type: object
                    type: array
                  commandOverrider:
                    description: CommandOverrider represents the rules dedicated to
                      handling container command
                    items:
                      description: CommandArgsOverrider represents the rules dedicated
                        to handling command/args overrides.
                      properties:
                        containerName:
                          description: The name of container
                          type: string
                        operator:
                          description: Operator represents the operator which will
                            apply on the command/args.
                          enum:
                          - add
                          - remove
                          type: string
                        value:
                          description: Value to be applied to command/args. Items
                            in Value which will be appended after command/args when
                            Operator is 'add'. Items in Value which match in command/args
                            will be deleted when Operator is 'remove'. If Value is
                            empty, then the command/args will remain the same.
                          items:
                            type: string
                          type: array
                      required:
                      - containerName
                      - operator
                      type: object
                    type: array
                  imageOverrider:
                    description: ImageOverrider represents the rules dedicated to
                      handling image overrides.
                    items:
                      description: ImageOverrider represents the rules dedicated to
                        handling image overrides.
                      properties:
                        component:
                          description: 'Component is part of image name. Basically
                            we presume an image can be made of ''[registry/]repository[:tag]''.
                            The registry could be: - k8s.gcr.io - fictional.registry.example:10443
                            The repository could be: - kube-apiserver - fictional/nginx
                            The tag cloud be: - latest - v1.19.1 - @sha256:dbcc1c35ac38df41fd2f5e4130b32ffdb93ebae8b3dbe638c23575912276fc9c'
                          enum:
                          - Registry
                          - Repository
                          - Tag
                          type: string
                        operator:
                          description: Operator represents the operator which will
                            apply on the image.
                          enum:
                          - add
                          - remove
                          - replace
                          type: string
                        predicate:
                          description: "Predicate filters images before applying the
                            rule. \n Defaults to nil, in that case, the system will
                            automatically detect image fields if the resource type
                            is Pod, ReplicaSet, Deployment or StatefulSet by following
                            rule:   - Pod: spec/containers/<N>/image   - ReplicaSet:
                            spec/template/spec/containers/<N>/image   - Deployment:
                            spec/template/spec/containers/<N>/image   - StatefulSet:
                            spec/template/spec/containers/<N>/image In addition, all
                            images will be processed if the resource object has more
                            than one containers. \n If not nil, only images matches
                            the filters will be processed."
                          properties:
                            path:
                              description: Path indicates the path of target field
                              type: string
                          required:
                          - path
                          type: object
                        value:
                          description: Value to be applied to image. Must not be empty
                            when operator is 'add' or 'replace'. Defaults to empty
                            and ignored when operator is 'remove'.
                          type: string
                      required:
                      - component
                      - operator
                      type: object
                    type: array
                  plaintext:
                    description: Plaintext represents override rules defined with
                      plaintext overriders.
                    items:
                      description: PlaintextOverrider is a simple overrider that overrides
                        target fields according to path, operator and value.
                      properties:
                        operator:
                          description: 'Operator indicates the operation on target
                            field. Available operators are: add, update and remove.'
                          enum:
                          - add
                          - remove
                          - replace
                          type: string
                        path:
                          description: Path indicates the path of target field
                          type: string
                        value:
                          description: Value to be applied to target field. Must be
                            empty when operator is Remove.
                          x-kubernetes-preserve-unknown-fields: true
                      required:
                      - operator
                      - path
                      type: object
                    type: array
                type: object
              resourceSelectors:
                description: ResourceSelectors restricts resource types that this
                  override policy applies to. nil means matching all resources.
                items:
                  description: ResourceSelector the resources will be selected.
                  properties:
                    apiVersion:
                      description: APIVersion represents the API version of the target
                        resources.
                      type: string
                    kind:
                      description: Kind represents the Kind of the target resources.
                      type: string
                    labelSelector:
                      description: A label query over a set of resources. If name
                        is not empty, labelSelector will be ignored.
                      properties:
                        matchExpressions:
                          description: matchExpressions is a list of label selector
                            requirements. The requirements are ANDed.
                          items:
                            description: A label selector requirement is a selector
                              that contains values, a key, and an operator that relates
                              the key and values.
                            properties:
                              key:
                                description: key is the label key that the selector
                                  applies to.
                                type: string
                              operator:
                                description: operator represents a key's relationship
                                  to a set of values. Valid operators are In, NotIn,
                                  Exists and DoesNotExist.
                                type: string
                              values:
                                description: values is an array of string values.
                                  If the operator is In or NotIn, the values array
                                  must be non-empty. If the operator is Exists or
                                  DoesNotExist, the values array must be empty. This
                                  array is replaced during a strategic merge patch.
                                items:
                                  type: string
                                type: array
                            required:
                            - key
                            - operator
                            type: object
                          type: array
                        matchLabels:
                          additionalProperties:
                            type: string
                          description: matchLabels is a map of {key,value} pairs.
                            A single {key,value} in the matchLabels map is equivalent
                            to an element of matchExpressions, whose key field is
                            "key", the operator is "In", and the values array contains
                            only "value". The requirements are ANDed.
                          type: object
                      type: object
                    name:
                      description: Name of the target resource. Default is empty,
                        which means selecting all resources.
                      type: string
                    namespace:
                      description: Namespace of the target resource. Default is empty,
                        which means inherit from the parent object scope.
                      type: string
                  required:
                  - apiVersion
                  - kind
                  type: object
                type: array
              targetCluster:
                description: TargetCluster defines restrictions on this override policy
                  that only applies to resources propagated to the matching clusters.
                  nil means matching all clusters.
                properties:
                  clusterNames:
                    description: ClusterNames is the list of clusters to be selected.
                    items:
                      type: string
                    type: array
                  exclude:
                    description: ExcludedClusters is the list of clusters to be ignored.
                    items:
                      type: string
                    type: array
                  fieldSelector:
                    description: FieldSelector is a filter to select member clusters
                      by fields. If non-nil and non-empty, only the clusters match
                      this filter will be selected.
                    properties:
                      matchExpressions:
                        description: A list of field selector requirements.
                        items:
                          description: A node selector requirement is a selector that
                            contains values, a key, and an operator that relates the
                            key and values.
                          properties:
                            key:
                              description: The label key that the selector applies
                                to.
                              type: string
                            operator:
                              description: Represents a key's relationship to a set
                                of values. Valid operators are In, NotIn, Exists,
                                DoesNotExist. Gt, and Lt.
                              type: string
                            values:
                              description: An array of string values. If the operator
                                is In or NotIn, the values array must be non-empty.
                                If the operator is Exists or DoesNotExist, the values
                                array must be empty. If the operator is Gt or Lt,
                                the values array must have a single element, which
                                will be interpreted as an integer. This array is replaced
                                during a strategic merge patch.
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                    type: object
                  labelSelector:
                    description: LabelSelector is a filter to select member clusters
                      by labels. If non-nil and non-empty, only the clusters match
                      this filter will be selected.
                    properties:
                      matchExpressions:
                        description: matchExpressions is a list of label selector
                          requirements. The requirements are ANDed.
                        items:
                          description: A label selector requirement is a selector
                            that contains values, a key, and an operator that relates
                            the key and values.
                          properties:
                            key:
                              description: key is the label key that the selector
                                applies to.
                              type: string
                            operator:
                              description: operator represents a key's relationship
                                to a set of values. Valid operators are In, NotIn,
                                Exists and DoesNotExist.
                              type: string
                            values:
                              description: values is an array of string values. If
                                the operator is In or NotIn, the values array must
                                be non-empty. If the operator is Exists or DoesNotExist,
                                the values array must be empty. This array is replaced
                                during a strategic merge patch.
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                      matchLabels:
                        additionalProperties:
                          type: string
                        description: matchLabels is a map of {key,value} pairs. A
                          single {key,value} in the matchLabels map is equivalent
                          to an element of matchExpressions, whose key field is "key",
                          the operator is "In", and the values array contains only
                          "value". The requirements are ANDed.
                        type: object
                    type: object
                type: object
            required:
            - overriders
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []`
