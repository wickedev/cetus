# Cetus

## 암묵적인 정보
- 아래 정보는 프로젝트의 .cetus/config에 저장된다.
- cluster namespace (user or group)
- image registry
- .cetus/cache 디렉토리 아래에 build에 관련된 캐시가 저장된다.

## 명시적인 정보
- 아래 정보는 \<project root\>/cetus.yaml에 개발자가 작성한다.
- images: build 명령어 실행시 빌드될 이미지 정보
- services: 배포될 서비스들(k8s deployment + service로 구성)
- dependencies: 의존성을 가지는 외부 서비스(git, path를 지원한다)
- profiles: patch(rfc6902), replace를 지원하여 
- dev: 개발시 필요한 설정들
- test: 배포 후 확인 커맨드 실행. 실패할 경우 기존 버전으로 rollback 된다.
- vars: file, env, stdin 입력을 받아 cetus.yaml 내에서 사용 할 수 있다.
- ingress: 외부에서 domain으로 서비스에 접근할 수 있다.

## 초기화
- cetus init
- cetus namespace add staging https://5c7110be.cetus.dev/wickedev/bookinfo
    - 배포시 클러스터 내에서는 wickedev--bookinfo k8s 네임스페이스를 가짐
- cetus namespace add prod https://5c7110be.cetus.dev/demo/bookinfo
    - 배포시 클러스터 내에서는 demo--bookinfo k8s 네임스페이스를 가짐

## 개발시
- cetus dev [cluster]
- [cluster] 인자 없이 실행 할 경우 로컬에 배포되며, 클러스터를 지정하여
- 해당 어플리케이션을 제외하고, 의존성이 있는 서비스(DB, MQ, 마이크로서비스 등)을 로컬 쿠버네티스 클러스에 배포한다.
- 배포된 각 서비스를 로컬 프로세스 0번 포트에 바인딩 하여 임의로 할당한 뒤, 현재 프로세스 환경 변수에 바인딩 한다.
    - 가령 postgresql 서비스가 로컬 30423 포트에 바인딩 되었다면, POSTGRESQL_SERVICE_HOST는 localhost, POSTGRESQL_SERVICE_PORT는 30423이다.
    - 해당 환경 변수는 어플리케이션을 실행하는 프로세스에만 바인딩 된다.

## 배포시
- cetus deploy staging
- cetus deploy prod
- deloy 호출 전에 publish(build, push), 호출 후에는 test가 불려진다.
- publish, build, push는 git 과 .cetus/cache 을 참조하여 캐시, 빌드, 버전을 결정한다.
    - 버전은 (branch)-(8 length git hash) 이다. 예시 master-1we09e41
    - 기본 동작은 git에 커밋된 코드만 빌드하지만 --force 플래그를 사용하면 커밋되지 않은 코드를 포함하며 버전 +(revision)이 붙는다. 예시 master-1we09e41+1
- cetus test
    - 

## 롤백
- cetus status [namespace]
    - k8s resources 및 그 상태 (uptime, age, restart, status)
- cetus history [namespace]
    - 지금까지 배포된 버전 목록
- cetus rollback [namespace] [version]

## 디버깅
- cetus debug [namespace] (pod) (-c [container])
    - 컨테이너가 죽었을 당시 pod log 및 events 출력
- cetus exec [namespace] (pod) (-c [container]) -- [exec]
    - 예) cetus exec staging -- ls -al
    - 컨테이너를 지정하지 않을 경우 대화식으로 선택 제공
- cetus attach [namespace] (pod) (-c [container])
    - 예) cetus attach staging
    - 컨테이너를 지정하지 않을 경우 대화식으로 선택 제공
- cetus logs (options) [namespace] (pod) (-c [container])
    - 예) cetus logs -f staging
    - 컨테이너를 지정하지 않을 경우 namespace의 전체 로그 (주입된 사이드카 제외)
- cetus curl [url]
    - 클러스터 URL(svc.cluster.local) 혹은 도메인으로 클러스터 내부에 curl을 수행
    - 예) cetus curl http://bookinfo.wickedev--bookinfo.svc.cluster.local
    - 예) cetus curl https://staging.bookinfo.5c7110be.cetus.dev

## 기타
- cetus ui
    - dev로 배포된 서비스 상태 확인 또는 관리 할 수 있는 UI를 브라우저에 띄운다.
