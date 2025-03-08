import {Node} from "matreshka-api/api/grpc/matreshka-be_api.pb.ts";
import {Telegram} from "@/models/AppConfig/Resources/Resource.ts";
import {extractStringValue} from "@/models/shared/common.ts";


export function mapTelegram(root: Node): Telegram {
    if (!root.name) {
        throw {message: "No data to parse telegram"};
    }
    const tg = new Telegram(root.name.slice(root.name.indexOf('TELEGRAM')).toLowerCase())


    root.innerNodes?.map(
        (n) => {
            if (!n.name || !root.name) {
                return
            }

            const fieldName = n.name.slice(root.name.length + 1)
            switch (fieldName) {
                case 'API-KEY':
                    tg.apiKey = extractStringValue(n)
            }
        }
    )

    return tg
}
